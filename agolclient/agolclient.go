package agolclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func SetConfig(cfg *Config) {
	config = cfg
}

var config = &Config{
	PortalAPIBaseUrl: "https://devext.arcgis.com/sharing/rest",
}
