package agolclient

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	TypeWebMap                = "Web Map"
	TypeWebMappingApplication = "Web Mapping Application"
	TypeKeywordRegisteredApp  = "Registered App"
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
	Id, Owner, Title, Type, Snippet, Thumbnail, FolderId string
	Modified                                             int64
	Tags                                                 []string
	TypeKeywords                                         []string
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

func (i *Item) HasTypeKeyword(typeKeyword string) bool {
	if i.TypeKeywords == nil {
		return false
	}
	for _, t := range i.TypeKeywords {
		if strings.EqualFold(typeKeyword, t) {
			return true
		}
	}
	return false
}

type RegisteredApp struct {
	ItemId, Client_Id, Client_Secret, AppType string
	Redirect_Uris                             []string
	Registered                                int64
}

func (r *RegisteredApp) RegisteredTime() time.Time {
	var t time.Time
	if r.Registered != 0 {
		t = time.Unix(0, r.Registered*int64(time.Millisecond))
	}
	return t
}

type RegisteredAppItem struct {
	*Item
	*RegisteredApp
}

func RegisteredAppItemsCsv(w io.Writer, ris []*RegisteredAppItem, portalHomeUrl string) {
	cw := csv.NewWriter(w)

	cw.Write([]string{"Title", "Item ID", "Client ID", "Client Secret", "Redirect URIs", "Registered", "Item URL"})
	for _, ri := range ris {
		cw.Write([]string{
			ri.Title,
			ri.Id,
			ri.Client_Id,
			ri.Client_Secret,
			strings.Join(ri.Redirect_Uris, ", "),
			ri.RegisteredTime().Format("January 1, 2006"),
			fmt.Sprintf("%s/item.html?id=%s", portalHomeUrl, ri.Id),
		})
	}

	cw.Flush()
}

type RegisteredAppLoginStatsReq struct {
	StartTime time.Time
	EndTime   time.Time
	Period    string
}

type RegisteredAppLoginStats struct {
	StartTime int64
	EndTime   int64
	Data      []*RegisteredAppLoginStatsData
}

func (r *RegisteredAppLoginStats) StartTimeTime() time.Time {
	var t time.Time
	if r.StartTime != 0 {
		t = time.Unix(0, r.StartTime*int64(time.Millisecond))
	}
	return t
}

func (r *RegisteredAppLoginStats) EndTimeTime() time.Time {
	var t time.Time
	if r.EndTime != 0 {
		t = time.Unix(0, r.EndTime*int64(time.Millisecond))
	}
	return t
}

func (r *RegisteredAppLoginStats) NumLogins() int {
	n := 0

	if r.Data != nil {
		for _, d := range r.Data {
			if d.Num != nil {
				for _, nums := range d.Num {
					num, _ := strconv.Atoi(nums[1])
					n = n + num
				}
			}
		}
	}

	return n
}

func (r *RegisteredAppLoginStats) NumOrgs() int {
	orgs := make(map[string]bool)

	if r.Data != nil {
		for _, d := range r.Data {
			if d.UserOrgId != "" {
				orgs[d.UserOrgId] = true
			}
		}
	}

	return len(orgs)
}

type RegisteredAppLoginStatsData struct {
	UserOrgId string
	Num       [][]string
}

type RegisteredAppItemLoginStats struct {
	*Item
	*RegisteredApp
	*RegisteredAppLoginStats
}

func RegisteredAppItemsLoginStatsCsv(w io.Writer, ris []*RegisteredAppItemLoginStats, portalHomeUrl string) {
	cw := csv.NewWriter(w)

	cw.Write([]string{"Application", "Users", "Logins", "Organizations", "Owner", "Item ID", "Item URL"})
	for _, ri := range ris {
		cw.Write([]string{
			ri.Title,
			strconv.Itoa(len(ri.Data)),
			strconv.Itoa(ri.NumLogins()),
			strconv.Itoa(ri.NumOrgs()),
			ri.Owner,
			ri.Id,
			fmt.Sprintf("%s/item.html?id=%s", portalHomeUrl, ri.Id),
		})
	}

	cw.Flush()
}

type RegisteredAppItemLoginStatsArray []*RegisteredAppItemLoginStats

func (r RegisteredAppItemLoginStatsArray) Len() int { return len(r) }

func (r RegisteredAppItemLoginStatsArray) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type RegisteredAppItemLoginStatsByNum struct {
	RegisteredAppItemLoginStatsArray
}

func (r *RegisteredAppItemLoginStatsByNum) Less(i, j int) bool {
	return len(r.RegisteredAppItemLoginStatsArray[i].Data) >= len(r.RegisteredAppItemLoginStatsArray[j].Data)
}

var _ sort.Interface = (*RegisteredAppItemLoginStatsByNum)(nil)

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
