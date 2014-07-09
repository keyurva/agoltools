package auth

import (
	"agoltools"
	"agoltools/config"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	agoltools.HandleFunc("/auth/signin", signIn)
	agoltools.HandleFunc("/auth/callback", callback)
	agoltools.HandleFunc("/auth/signout", signOut)

	authParams := url.Values{
		"client_id":     {config.Config.AppId},
		"redirect_uri":  {config.Config.AppBaseUrl + "/auth/callback"},
		"expiration":    {strconv.Itoa(60 * 24 * 14)}, // 2 weeks in minutes
		"response_type": {"token"},
	}
	portalAuthorizeUrl = fmt.Sprintf("%s/oauth2/authorize?%s", config.Config.PortalAPIBaseUrl, authParams.Encode())
	portalOrgAuthorizeUrlTemplate = fmt.Sprintf("%s/oauth2/authorize?%s", config.Config.PortalOrgAPIBaseUrlTemplate, authParams.Encode())
}

func Authenticated(handler agoltools.HandlerFunc) (authenticated agoltools.HandlerFunc) {
	return func(r *agoltools.Request) (err error) {
		if r.Auth == nil {
			r.Redirect("/auth/signin?redirect=" + url.QueryEscape(r.URLWithQuery()))
			return nil
		}
		return handler(r)
	}
}

var portalAuthorizeUrl, portalOrgAuthorizeUrlTemplate string

const (
	signInTemplate   = "agoltools/auth/templates/signin.html"
	callbackTemplate = "agoltools/auth/templates/callback.html"
)

func signIn(r *agoltools.Request) (err error) {
	authUrl := portalAuthorizeUrl

	if redirect := r.R.FormValue("redirect"); strings.HasPrefix(redirect, "/") {
		redirect = url.QueryEscape(redirect)
		authUrl += "&state=" + redirect
	}

	r.Redirect(authUrl)

	return nil
}

func callback(r *agoltools.Request) (err error) {
	r.AddData(map[string]interface{}{
		"PageTitle":      "Sign in to ArcGIS Online",
		"AuthCookieName": agoltools.AuthCookie,
	})

	return r.RenderUsingBaseTemplate(callbackTemplate)
}

func signOut(r *agoltools.Request) (err error) {
	r.SignOut()
	r.Redirect("/")
	return nil
}
