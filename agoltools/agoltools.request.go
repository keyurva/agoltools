package agoltools

import (
	"agolclient"
	"agoltools/config"
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

var timeoutDuration = time.Second * 60

type Request struct {
	R         *http.Request
	W         http.ResponseWriter
	Auth      *agolclient.Auth
	Data      map[string]interface{}
	Context   appengine.Context
	Cache     *Cache
	transport *urlfetch.Transport
	self      *agolclient.PortalSelf
	selfChan  chan *agolclient.PortalSelf
}

func (r *Request) PortalSelf() *agolclient.PortalSelf {
	if r.Auth == nil {
		return nil
	}
	if r.self != nil {
		return r.self
	}
	if r.selfChan != nil {
		r.self = <-r.selfChan
	}
	return r.self
}

func (r *Request) Org() *agolclient.Org {
	self := r.PortalSelf()
	if self == nil {
		return nil
	}
	return self.Org
}

func (r *Request) OrgUrlKey() string {
	org := r.Org()
	if org == nil {
		return ""
	}
	return org.UrlKey
}

func (r *Request) OrgId() string {
	org := r.Org()
	if org == nil {
		return ""
	}
	return org.Id
}

func (r *Request) PortalHomeUrl() string {
	return config.PortalHomeUrl(r.OrgUrlKey())
}

func (r *Request) RenderTemplates(templates ...string) (err error) {
	tname := filepath.Base(templates[0])
	t, err := template.New(tname).Funcs(templateFuncs).ParseFiles(templates...)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(r.W, tname, r); err != nil {
		return err
	}

	return nil
}

func (r *Request) LogInfof(format string, args ...interface{}) {
	r.Context.Infof(format, args...)
}

func (r *Request) LogErrorf(format string, args ...interface{}) {
	r.Context.Errorf(format, args...)
}

func (r *Request) Transport() *urlfetch.Transport {
	if r.transport == nil {
		r.transport = &urlfetch.Transport{
			Context:                       r.Context,
			Deadline:                      timeoutDuration,
			AllowInvalidServerCertificate: true,
		}
	}
	return r.transport
}

func (r *Request) URLWithQuery() string {
	u := r.R.URL.Path
	if q := r.R.URL.RawQuery; len(q) > 0 {
		u += "?" + q
	}
	return u
}

func (r *Request) SignOut() {
	http.SetCookie(r.W, &http.Cookie{Name: AuthCookie, Value: "", Path: "/", MaxAge: -1})
}

func (r *Request) Redirect(url string) {
	http.Redirect(r.W, r.R, url, http.StatusTemporaryRedirect)
}

func (r *Request) AddData(data map[string]interface{}) {
	for k, v := range data {
		r.Data[k] = v
	}
}

func (r *Request) getPortalSelfAsync() {
	if r.Auth != nil {
		// kick off a concurrent request to get portal self
		r.selfChan = make(chan *agolclient.PortalSelf, 1)
		go func() {
			selfKey := fmt.Sprintf("self_%s", r.Auth.Username)
			var self agolclient.PortalSelf
			err := r.Cache.Get(selfKey, &self)
			if err == nil {
				r.selfChan <- &self
				return
			}

			r.self, err = agolclient.GetPortalSelf(r.Transport(), r.Auth)
			if err != nil {
				r.selfChan <- nil
			} else {
				r.selfChan <- r.self
				r.Cache.SetWithExpiration(selfKey, r.self, time.Minute*10)
			}
		}()
	}
}

func newRequest(w http.ResponseWriter, r *http.Request) (req *Request) {
	auth, _ := getAuthFromCookie(r)
	c := appengine.NewContext(r)
	req = &Request{
		R:       r,
		W:       w,
		Auth:    auth,
		Data:    make(map[string]interface{}),
		Context: c,
		Cache:   NewCache(c),
	}

	req.getPortalSelfAsync()

	return req
}

func getAuthFromCookie(r *http.Request) (auth *agolclient.Auth, err error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return nil, err
	}
	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(value), &auth); err != nil {
		return nil, err
	}
	return auth, nil
}
