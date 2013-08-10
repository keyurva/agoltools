package agolclient

import (
	"fmt"
	"net/url"
	"time"
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
