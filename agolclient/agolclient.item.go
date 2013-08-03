package agolclient

import (
	"net/url"
)

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
