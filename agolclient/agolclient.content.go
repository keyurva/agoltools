package agolclient

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
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

func (w *WebMap) NumLayers() int {
	return w.NumOperationalLayers() + w.NumBaseMapLayers()
}

func (w *WebMap) NumOperationalLayers() int {
	num := 0
	for _, layer := range w.OperationalLayers {
		if len(layer.Url) > 0 {
			num++
		}
	}
	return num
}

func (w *WebMap) NumBaseMapLayers() int {
	num := 0
	for _, layer := range w.BaseMap.BaseMapLayers {
		if len(layer.Url) > 0 {
			num++
		}
	}
	return num
}

type WebMapItem struct {
	*Item
	*WebMap
}

func WebMapItemsCsv(w io.Writer, wmis []*WebMapItem, portalHomeUrl string) {
	cw := csv.NewWriter(w)

	cw.Write([]string{"Title", "Owner", "Item ID", "Last Modified", "Item URL"})
	for _, wmi := range wmis {
		cw.Write([]string{
			wmi.Title,
			wmi.Owner,
			wmi.Id,
			wmi.ModifiedTime().Format("January 1, 2006"),
			fmt.Sprintf("%s/item.html?id=%s", portalHomeUrl, wmi.Id),
		})
	}

	cw.Flush()
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
