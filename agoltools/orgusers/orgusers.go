package orgusers

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"strings"
)

func init() {
	agoltools.HandleFunc("/orgusers", auth.Authenticated(getOrgUsers))
}

const (
	getOrgUsersTemplate = "agoltools/orgusers/templates/getorgusers.html"
	orgUsersTemplate    = "agoltools/orgusers/templates/orgusers.html"
)

func getOrgUsers(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "Get Organization Users"

	f := strings.ToLower(strings.Trim(r.R.FormValue("f"), " "))
	if f == "" {
		return r.RenderUsingBaseTemplate(getOrgUsersTemplate)
	}

	us, err := agolclient.GetAllOrgUsers(r.Transport(), r.Auth)
	if err != nil {
		return err
	}

	if f == "csv" {
		r.W.Header().Set("Content-Type", "text/csv")
		r.W.Header().Set("Content-Disposition", "inline;filename=orgusers.csv")
		agolclient.UsersCsv(r.W, us)
		return
	}

	r.Data["Users"] = us

	return r.RenderUsingBaseTemplate(orgUsersTemplate)
}
