package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"net/http"
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
	r.Data["PageTitle"] = "Find Web Maps With URL"

	url := strings.ToLower(strings.Trim(r.R.FormValue("url"), " "))
	if url == "" {
		return r.RenderUsingBaseTemplate(getWebMapUrlsTemplate)
	}

	findFor := strings.ToLower(strings.Trim(r.R.FormValue("for"), " "))

	var wmis []*agolclient.WebMapItem

	if findFor == "org" {
		accountId := r.OrgId()
		if accountId == "" {
			return &agoltools.Error{
				Message: "This option is only available to users that belong to an organization",
				Code:    http.StatusBadRequest,
			}
		}

		wmis, err = agolclient.GetOrgWebMapsWithUrl(r.Transport(), accountId, url, r.Auth)
		if err != nil {
			return err
		}
	} else {
		wmis, err = agolclient.GetUserWebMapsWithUrl(r.Transport(), url, r.Auth)
		if err != nil {
			return err
		}
	}

	f := strings.ToLower(strings.Trim(r.R.FormValue("f"), " "))
	if f == "csv" {
		r.W.Header().Set("Content-Type", "text/csv")
		r.W.Header().Set("Content-Disposition", "inline;filename=webmaps.csv")
		agolclient.WebMapItemsCsv(r.W, wmis, r.PortalHomeUrl())
		return
	}

	r.AddData(map[string]interface{}{
		"PageTitle":   "Web Maps With URL",
		"WebMapItems": wmis,
		"URL":         strings.Trim(r.R.FormValue("url"), " "),
		"For":         findFor,
	})

	return r.RenderUsingBaseTemplate(webMapUrlsTemplate)
}
