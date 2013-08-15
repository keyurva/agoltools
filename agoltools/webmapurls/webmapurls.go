package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"strings"
)

func init() {
	agoltools.HandleFunc("/webmapurls", auth.Authenticated(getWebMapUrls))
}

const (
	getWebMapUrlsTemplate = "agoltools/webmapurls/templates/getwebmapurls.html"
	webMapUrlsTemplate    = "agoltools/webmapurls/templates/webmapurls.html"
)

func getWebMapUrls(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "Get Web Maps With URL"

	url := strings.ToLower(strings.Trim(r.R.FormValue("url"), " "))
	if url == "" {
		return r.RenderUsingBaseTemplate(getWebMapUrlsTemplate)
	}

	wmis, err := agolclient.GetUserWebMapsWithUrl(r.Transport(), url, r.Auth)
	if err != nil {
		return err
	}

	r.AddData(map[string]interface{}{
		"PageTitle":   "Web Maps With URL",
		"WebMapItems": wmis,
		"URL":         r.R.FormValue("url"),
	})

	return r.RenderUsingBaseTemplate(webMapUrlsTemplate)
}
