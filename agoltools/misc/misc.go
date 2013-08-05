package misc

import (
	"agoltools"
)

const (
	aboutTemplate = "agoltools/misc/templates/about.html"
)

func init() {
	agoltools.HandleFunc("/about", root)
}

func root(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "About ArcGIS Online Tools"
	return r.RenderUsingBaseTemplate(aboutTemplate)
}
