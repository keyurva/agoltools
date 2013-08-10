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

	// panel dropdown and ids
	pdropdown := []string{}
	pids := map[string]string{} //[display name]id

	addPanel := func(condition bool, name, id string) {
		if condition {
			pdropdown = append(pdropdown, name)
			pids[name] = id
		}
	}

	addPanel(myagol.User != nil, "User Info", "user-panel")
	addPanel(myagol.Folders != nil, "My Content", "content-panel")
	addPanel(myagol.User != nil && myagol.User.Groups != nil, "My Groups", "groups-panel")
	addPanel(myagol.Org != nil, "Organization Info", "org-panel")
	addPanel(myagol.Subscription != nil, "Subscription Info", "sub-panel")

	if len(pdropdown) > 1 {
		r.Data["PanelDropdown"] = pdropdown
		r.Data["PanelIds"] = pids
	}

	return r.RenderUsingBaseTemplate(orgUsersTemplate)
}
