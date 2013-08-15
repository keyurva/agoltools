package agolclient

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	TypeWebMap                = "Web Map"
	TypeWebMappingApplication = "Web Mapping Application"
)

type Folder struct {
	Username, Id, Title string
}

type FolderContent struct {
	Folder
	Items   []*Item
	Folders []*Folder
}

type Item struct {
	Id, Owner, Title, Type, Snippet, Thumbnail string
	Modified                                   int64
	Tags                                       []string
}

func (i *Item) ModifiedTime() *time.Time {
	if i.Modified == 0 {
		return nil
	}
	t := time.Unix(0, i.Modified*int64(time.Millisecond))
	return &t
}

func (i *Item) RelativeThumbnailUrl() string {
	if i.Thumbnail == "" {
		return ""
	}
	return fmt.Sprintf("/content/items/%s/info/%s", i.Id, i.Thumbnail)
}

type WebMap struct {
	OperationalLayers []struct {
		Url, Id, Title, ItemId string
	}
	BaseMap struct {
		BaseMapLayers []struct {
			Id, Url string
		}
	}
}

func (w *WebMap) HasUrl(url string) bool {
	for _, l := range w.OperationalLayers {
		if strings.EqualFold(url, l.Url) {
			return true
		}
	}
	for _, l := range w.BaseMap.BaseMapLayers {
		if strings.EqualFold(url, l.Url) {
			return true
		}
	}
	return false
}

type WebMapItem struct {
	*Item
	*WebMap
}

type SearchResponse struct {
	Total     int
	Start     int
	Num       int
	NextStart int
	Results   []*Item
}

func (sr *SearchResponse) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Total: %d, Start: %d, Num %d, Next Start: %d", sr.Total, sr.Start, sr.Num, sr.NextStart)

	for _, i := range sr.Results {
		fmt.Fprintf(&buf, "\n%v", i)
	}

	return buf.String()
}

type ItemMapper interface {
	ItemMap() map[string]string
}

func addItemValues(im ItemMapper, values url.Values) {
	for k, v := range im.ItemMap() {
		values.Add(k, v)
	}
}

type AddItemResponse struct {
	Success    bool
	Id, Folder string
}
