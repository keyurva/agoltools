package addfolderservices

import (
	"agolclient"
	"agoltools"
	"agoltools/auth"
	"strings"
)

func init() {
	agoltools.HandleFunc("/addfolderservices", auth.Authenticated(addFolderServices))
}

const (
	addFolderServicesTemplate = "agoltools/addfolderservices/templates/addfolderservices.html"
	folderServicesTemplate    = "agoltools/addfolderservices/templates/folderservices.html"
)

func addFolderServices(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "Add Folder Services"

	folderUrl := strings.Trim(r.R.FormValue("folderUrl"), " ")
	if folderUrl == "" {
		return r.RenderUsingBaseTemplate(addFolderServicesTemplate)
	}

	folder, catalog, status, err := agolclient.AddFolderServices(r.Transport(), folderUrl, r.Auth)
	if err != nil {
		return err
	}

	r.Data["Folder"] = folder
	r.Data["Catalog"] = catalog
	r.Data["Status"] = status

	return r.RenderUsingBaseTemplate(folderServicesTemplate)
}
