package root

import (
	"agoltools"
)

const (
	rootTemplate = "agoltools/root/templates/root.html"
)

func init() {
	agoltools.HandleFunc("/", root)
}

func root(r *agoltools.Request) (err error) {
	r.Data["PageTitle"] = "ArcGIS Online Tools"
	return r.RenderUsingBaseTemplate(rootTemplate)
}
