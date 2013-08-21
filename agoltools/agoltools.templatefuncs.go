package agoltools

import (
	"agolclient"
	"agoltools/config"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

var templateFuncs = template.FuncMap{
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
	"gt": func(first, second interface{}) bool {
		switch first.(type) {
		case int, int8, int16, int32, int64:
			switch second.(type) {
			case int, int8, int16, int32, int64:
				return reflect.ValueOf(first).Int() > reflect.ValueOf(second).Int()
			}
		}
		return false
	},
	"eq": func(first, second interface{}) bool {
		switch first.(type) {
		case int, int8, int16, int32, int64:
			switch second.(type) {
			case int, int8, int16, int32, int64:
				return reflect.ValueOf(first).Int() == reflect.ValueOf(second).Int()
			}
		case string:
			switch second.(type) {
			case string:
				return strings.EqualFold(reflect.ValueOf(first).String(), reflect.ValueOf(second).String())
			}
		}
		return false
	},
	"portalUrl": func(relativeUrl string, auth *agolclient.Auth) string {
		portalUrl := fmt.Sprintf("%s%s", config.Config.PortalAPIBaseUrl, relativeUrl)
		if auth != nil {
			portalUrl = fmt.Sprintf("%s?token=%s", portalUrl, auth.AccessToken)
		}
		return portalUrl
	},
}

func (r *Request) RenderUsingBaseTemplate(templateFilePaths ...string) (err error) {
	templates := []string{baseTemplate, headerTemplate}
	templates = append(templates, templateFilePaths...)
	return r.RenderTemplates(templates...)
}
