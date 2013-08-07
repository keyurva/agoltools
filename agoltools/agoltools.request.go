package agoltools

import (
	"agolclient"
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	R         *http.Request
	W         http.ResponseWriter
	Auth      *agolclient.Auth
	Data      map[string]interface{}
	context   appengine.Context
	transport *urlfetch.Transport
}

var timeoutDuration = time.Second * 60

var templateFuncs = template.FuncMap{
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func (r *Request) RenderUsingBaseTemplate(templateFilePaths ...string) (err error) {
	templates := []string{baseTemplate, headerTemplate}
	templates = append(templates, templateFilePaths...)
	return r.RenderTemplates(templates...)
}

func (r *Request) RenderTemplates(templates ...string) (err error) {
	t := template.New("").Funcs(templateFuncs)
	// reading files instead of calling ParseFiles directly because
	// I haven't been able to get Funcs to work when using ParseFiles so far
	for _, tmpl := range templates {
		b, err := ioutil.ReadFile(tmpl)
		if err != nil {
			return err
		}
		t, err = t.Parse(string(b))
		if err != nil {
			return err
		}
	}

	if err := t.Execute(r.W, r); err != nil {
		return err
	}

	return nil
}

func (r *Request) Context() appengine.Context {
	if r.context == nil {
		r.context = appengine.NewContext(r.R)
	}
	return r.context
}

func (r *Request) LogInfof(format string, args ...interface{}) {
	r.Context().Infof(format, args...)
}

func (r *Request) Transport() *urlfetch.Transport {
	if r.transport == nil {
		r.transport = &urlfetch.Transport{
			Context:                       r.Context(),
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

func newRequest(w http.ResponseWriter, r *http.Request) (req *Request) {
	auth, _ := getAuthFromCookie(r)
	return &Request{
		R:    r,
		W:    w,
		Auth: auth,
		Data: make(map[string]interface{}),
	}
}

func getAuthFromCookie(r *http.Request) (auth *agolclient.Auth, err error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return
	}
	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(value), &auth); err != nil {
		return
	}
	return auth, nil
}
