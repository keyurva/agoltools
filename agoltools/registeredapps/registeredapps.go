package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"strings"
)

func init() {
	agoltools.HandleFunc("/registeredapps", auth.Authenticated(getRegisteredApps))
}

const (
	getRegisteredAppsTemplate = "agoltools/registeredapps/templates/getregisteredapps.html"
	registeredAppsTemplate    = "agoltools/registeredapps/templates/registeredapps.html"
)

func getRegisteredApps(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "My Registered Apps"

	f := strings.ToLower(strings.Trim(r.R.FormValue("f"), " "))
	if f == "" {
		return r.RenderUsingBaseTemplate(getRegisteredAppsTemplate)
	}

	ris, err := agolclient.GetUserRegisteredApps(r.Transport(), r.Auth)
	if err != nil {
		return err
	}

	if f == "csv" {
		r.W.Header().Set("Content-Type", "text/csv")
		r.W.Header().Set("Content-Disposition", "inline;filename=registeredapps.csv")
		agolclient.RegisteredAppItemsCsv(r.W, ris, r.PortalHomeUrl())
		return
	}

	r.AddData(map[string]interface{}{
		"RegisteredAppItems": ris,
	})

	return r.RenderUsingBaseTemplate(registeredAppsTemplate)
}
