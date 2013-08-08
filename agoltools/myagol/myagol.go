package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"agoltools/config"
	"strings"
)

func init() {
	agoltools.HandleFunc("/myagol", auth.Authenticated(getMyAGOL))
}

const (
	getOrgUsersTemplate = "agoltools/myagol/templates/getmyagol.html"
	orgUsersTemplate    = "agoltools/myagol/templates/myagol.html"
)

func getMyAGOL(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "Get My ArcGIS Online Information"

	f := strings.ToLower(strings.Trim(r.R.FormValue("f"), " "))
	if f == "" {
		return r.RenderUsingBaseTemplate(getOrgUsersTemplate)
	}

	myagol, err := agolclient.GetMyAGOL(r.Transport(), r.Auth)
	if err != nil {
		return err
	}

	r.Data["PageTitle"] = "My ArcGIS Online Information"

	r.Data["MyAGOL"] = myagol

	orgUrlKey := ""
	if myagol.Org != nil {
		orgUrlKey = myagol.Org.UrlKey
	}
	r.Data["PortalHomeUrl"] = config.PortalHomeUrl(orgUrlKey)

	return r.RenderUsingBaseTemplate(orgUsersTemplate)
}
