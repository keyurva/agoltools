package agolclient

import (
	"bytes"
	"fmt"
)

var ServiceTypes = map[string]string{
	"GeocodeServer":  "Geocoding Service",
	"GeoDataServer":  "Geodata Service",
	"GeometryServer": "Geometry Service",
	"GlobeServer":    "Globe Service",
	"GPServer":       "Geoprocessing Service",
	"ImageServer":    "Image Service",
	"FeatureServer":  "Feature Service",
	"MapServer":      "Map Service",
	"NAServer":       "Network Analysis Service",
}

type CatalogService struct {
	Name string
	Type string
	Url  string
}

func (cs *CatalogService) String() string {
	return fmt.Sprintf("%s, %s, %s", cs.Name, cs.Type, cs.Url)
}

func (cs *CatalogService) ItemMap() map[string]string {
	m := make(map[string]string)
	m["title"] = cs.Name + "/" + cs.Type
	m["type"] = ServiceTypes[cs.Type]
	m["url"] = cs.Url
	return m
}

type ServiceCatalog struct {
	Services []*CatalogService
}

func (sc *ServiceCatalog) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Services (%d):", len(sc.Services)))

	for _, s := range sc.Services {
		buf.WriteString(fmt.Sprintf("\n\t(%s)", s))
	}

	return buf.String()
}
