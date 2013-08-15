package agolclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// arbitrary limit on max search items returned
	MaxSearchItems = 10000
)

func GetOrgWebMapsWithUrl(rt http.RoundTripper, accountId string, url string, auth *Auth) (wis []*WebMapItem, err error) {
	items, err := GetAllSearchItems(rt, fmt.Sprintf(`accountid:%s type:"%s" -type:"%s"`, accountId, TypeWebMap, TypeWebMappingApplication), auth)
	if err != nil {
		return nil, DisplayError("Unable to get organization webmaps", err)
	}
	return GetWebMapsWithUrl(rt, items, url, auth)
}

func GetUserWebMapsWithUrl(rt http.RoundTripper, url string, auth *Auth) (wis []*WebMapItem, err error) {
	items, err := GetAllSearchItems(rt, fmt.Sprintf(`owner:%s type:"%s" -type:"%s"`, auth.Username, TypeWebMap, TypeWebMappingApplication), auth)
	if err != nil {
		return nil, DisplayError("Unable to get user webmaps", err)
	}
	return GetWebMapsWithUrl(rt, items, url, auth)
}

func GetWebMapsWithUrl(rt http.RoundTripper, items []*Item, url string, auth *Auth) (wis []*WebMapItem, err error) {
	wiChan := make(chan *WebMapItem)
	num := 0

	for _, i := range items {
		if i.Type == TypeWebMap {
			wi := &WebMapItem{}
			wi.Item = i
			num++
			go func(wi *WebMapItem) {
				var wm WebMap
				err := GetItemData(rt, wi.Item.Id, &wm, auth)
				if err == nil && (&wm).HasUrl(url) {
					wi.WebMap = &wm
					wiChan <- wi
				} else {
					wiChan <- nil
				}

			}(wi)
		}
	}

	wis = []*WebMapItem{}
	for i := 0; i < num; i++ {
		wi := <-wiChan
		if wi != nil {
			wis = append(wis, wi)
		}
	}

	return wis, nil
}

func GetAllSearchItems(rt http.RoundTripper, q string, auth *Auth) (items []*Item, err error) {
	items = []*Item{}
	// get first batch
	start := 1
	num := 100
	sr, err := SearchItems(rt, q, start, num, auth)
	if err != nil {
		return nil, DisplayError("Unable to get search items", err)
	}
	items = append(items, sr.Results...)

	// concurrently fetch other batches based on total
	total := sr.Total
	start += num
	batches := make(chan []*Item)
	numBatches := 0
	for start <= total && start <= MaxSearchItems {
		go func(start int) {
			sr, _ := SearchItems(rt, q, start, num, auth)
			if sr != nil {
				batches <- sr.Results
			} else {
				batches <- nil
			}
		}(start)
		numBatches++
		start += num
	}

	for i := 0; i < numBatches; i++ {
		is := <-batches
		if is != nil {
			items = append(items, is...)
		}
	}

	return items, nil
}

func SearchItems(rt http.RoundTripper, q string, start int, num int, auth *Auth) (sr *SearchResponse, err error) {
	params := url.Values{"f": {"json"}, "q": {q}, "start": {strconv.Itoa(start)}, "num": {strconv.Itoa(num)}}
	if auth != nil {
		params.Add("token", auth.AccessToken)
	}
	url := fmt.Sprintf("%s/search", config.PortalAPIBaseUrl)

	if err = getAndUnmarshalJson(rt, url, params, &sr); err != nil {
		return nil, err
	}

	return sr, nil
}

func GetItemData(rt http.RoundTripper, itemId string, data interface{}, auth *Auth) (err error) {
	params := url.Values{"f": {"json"}}
	if auth != nil {
		params.Add("token", auth.AccessToken)
	}
	url := fmt.Sprintf("%s/content/items/%s/data", config.PortalAPIBaseUrl, itemId)

	if err = getAndUnmarshalJson(rt, url, params, &data); err != nil {
		return err
	}

	return nil
}

func GetMyAGOL(rt http.RoundTripper, auth *Auth) (ma *MyAGOL, err error) {
	selfChan := make(chan *PortalSelf)
	userChan := make(chan *User)
	contentChan := make(chan []*FolderContent)

	go func() {
		s, err := GetPortalSelf(rt, auth)
		if err != nil {
			selfChan <- nil
		} else {
			selfChan <- s
		}
	}()

	go func() {
		u, err := GetUser(rt, auth)
		if err != nil {
			userChan <- nil
		} else {
			userChan <- u
		}
	}()

	go func() {
		fs, err := GetUserContent(rt, auth)
		if err != nil {
			contentChan <- nil
		} else {
			contentChan <- fs
		}
	}()

	ma = &MyAGOL{}

	s := <-selfChan
	if s == nil {
		return nil, DisplayError("Unable to get self response", nil)
	}

	if s.Id != "" {
		ma.Org = s.Org
		ma.Subscription = s.SubscriptionInfo
	}

	u := <-userChan
	if s == nil {
		return nil, DisplayError("Unable to get user", nil)
	}

	ma.User = u

	fs := <-contentChan
	if fs == nil {
		return nil, DisplayError("Unable to get user content", nil)
	}

	ma.Folders = fs

	numItems := 0
	for _, f := range fs {
		if f.Items != nil {
			numItems += len(f.Items)
		}
	}
	ma.NumItems = numItems

	return ma, nil
}

func GetUserContent(rt http.RoundTripper, auth *Auth) (fs []*FolderContent, err error) {
	fs = []*FolderContent{}
	// get root folder
	root, err := GetFolderContent(rt, "", auth)
	if err != nil {
		return nil, err
	}

	fs = append(fs, root)

	// get all subfolders concurrently
	if root.Folders != nil {
		fchan := make(chan *FolderContent)
		for _, f := range root.Folders {
			go func(f *Folder) {
				fc, err := GetFolderContent(rt, f.Id, auth)
				if err != nil {
					LogError(err, true)
					fchan <- nil
				} else {
					fc.Folder = *f
					fchan <- fc
				}
			}(f)
		}

		for _, _ = range root.Folders {
			fc := <-fchan
			if fc != nil {
				fs = append(fs, fc)
			}
		}
	}

	return fs, nil
}

func GetFolderContent(rt http.RoundTripper, folderId string, auth *Auth) (f *FolderContent, err error) {
	folderUri := folderId
	if folderUri != "" {
		folderUri = "/" + folderUri
	}
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}}
	url := fmt.Sprintf("%s/content/users/%s%s", config.PortalAPIBaseUrl, auth.Username, folderUri)

	if err = getAndUnmarshalJson(rt, url, params, &f); err != nil {
		return nil, err
	}

	f.Id = folderId

	return f, nil
}

func GetUser(rt http.RoundTripper, auth *Auth) (u *User, err error) {
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}}
	url := fmt.Sprintf("%s/community/users/%s", config.PortalAPIBaseUrl, auth.Username)

	if err = getAndUnmarshalJson(rt, url, params, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func GetPortalSelf(rt http.RoundTripper, auth *Auth) (self *PortalSelf, err error) {
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}}
	url := fmt.Sprintf("%s/portals/self", config.PortalAPIBaseUrl)

	if err = getAndUnmarshalJson(rt, url, params, &self); err != nil {
		LogError(err, true)
		return nil, err
	}

	return self, nil
}

func GetAllOrgUsers(rt http.RoundTripper, auth *Auth) (users []*User, err error) {
	users = []*User{}
	// get first batch
	start := 1
	num := 100
	ur, err := GetOrgUsers(rt, start, num, auth)
	if err != nil {
		return nil, DisplayError("Unable to get organization users", err)
	}
	users = append(users, ur.Users...)

	// concurrently fetch other batches based on total
	total := ur.Total
	start += num
	batches := make(chan []*User)
	numBatches := 0
	for start <= total {
		go func(start int) {
			ur, _ := GetOrgUsers(rt, start, num, auth)
			if ur != nil {
				batches <- ur.Users
			} else {
				batches <- nil
			}
		}(start)
		numBatches++
		start += num
	}

	for i := 0; i < numBatches; i++ {
		us := <-batches
		if us != nil {
			users = append(users, us...)
		}
	}

	return users, nil
}

func GetOrgUsers(rt http.RoundTripper, start int, num int, auth *Auth) (ur *UsersResponse, err error) {
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}, "start": {strconv.Itoa(start)}, "num": {strconv.Itoa(num)}}
	url := fmt.Sprintf("%s/portals/self/users", config.PortalAPIBaseUrl)

	if err = getAndUnmarshalJson(rt, url, params, &ur); err != nil {
		return nil, err
	}

	return ur, nil
}

func AddFolderServices(rt http.RoundTripper, folderUrl string, auth *Auth) (folder *Folder, catalog *ServiceCatalog, status map[string]bool, err error) {
	// get catalog
	catalog, err = GetServiceCatalog(rt, folderUrl)
	if err != nil {
		return nil, catalog, nil, DisplayError("Unable to get folder services. Please check the folder URL you provided.", err)
	}

	// create folder
	folderTitle := fmt.Sprintf("%s - %s", folderUrl[strings.LastIndex(folderUrl, "/")+1:], randString(10))
	folder, err = CreateFolder(rt, folderTitle, auth)
	if err != nil {
		return folder, catalog, nil, DisplayError("Unable to create folder", err)
	}

	// add service items concurrently
	type Status struct {
		url     string
		success bool
	}

	done := make(chan *Status)
	for _, service := range catalog.Services {
		go func(service *CatalogService) {
			_, err = AddItem(rt, service, folder.Id, auth)
			success := true
			if err != nil {
				success = false
			}
			done <- &Status{service.Url, success}
		}(service)
	}

	status = make(map[string]bool)

	for _, _ = range catalog.Services {
		st := <-done
		status[st.url] = st.success
	}

	// return
	return folder, catalog, status, nil
}

func AddItem(rt http.RoundTripper, item ItemMapper, folderId string, auth *Auth) (res *AddItemResponse, err error) {
	if folderId != "" {
		folderId += "/"
	}
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}}
	addItemValues(item, params)
	url := fmt.Sprintf("%s/content/users/%s/%saddItem", config.PortalAPIBaseUrl, auth.Username, folderId)

	if err := postAndUnmarshalJson(rt, url, params, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func CreateFolder(rt http.RoundTripper, folderTitle string, auth *Auth) (folder *Folder, err error) {
	params := url.Values{"f": {"json"}, "token": {auth.AccessToken}, "title": {folderTitle}}
	url := fmt.Sprintf("%s/content/users/%s/createFolder", config.PortalAPIBaseUrl, auth.Username)

	var res = struct {
		Folder *Folder
	}{}

	if err := postAndUnmarshalJson(rt, url, params, &res); err != nil {
		return nil, err
	}

	return res.Folder, nil
}

func GetServiceCatalog(rt http.RoundTripper, folderUrl string) (catalog *ServiceCatalog, err error) {
	params := url.Values{"f": {"json"}}
	if err = getAndUnmarshalJson(rt, folderUrl, params, &catalog); err != nil {
		return nil, err
	}
	for _, s := range catalog.Services {
		s.Url = strings.Join([]string{folderUrl, s.Name[(strings.LastIndex(s.Name, "/") + 1):], s.Type}, "/")
	}
	return catalog, nil
}

func GenerateToken(rt http.RoundTripper, username string, password string) (auth *Auth, err error) {
	params := url.Values{"f": {"json"}, "username": {username}, "password": {password}, "client": {"requestip"}, "expiration": {"20160"}}
	url := fmt.Sprintf("%s/generateToken", config.PortalAPIBaseUrl)

	type Response struct{ Token string }

	var res Response

	if err = postAndUnmarshalJson(rt, url, params, &res); err != nil {
		return nil, err
	}

	return &Auth{Username: username, AccessToken: res.Token}, nil

}

func SetConfig(cfg *Config) {
	config = cfg
}

var (
	DevExtConfig = &Config{
		PortalAPIBaseUrl: "https://devext.arcgis.com/sharing/rest",
	}
	ProdConfig = &Config{
		PortalAPIBaseUrl: "https://www.arcgis.com/sharing/rest",
	}
	config = DevExtConfig
)
