package agoltools

import (
	"agolclient"
	"net/http"
)

const (
	AuthCookie                 = "agoltools_auth"
	baseTemplate               = "agoltools/templates/base.html"
	headerTemplate             = "agoltools/templates/header.html"
	unsupportedBrowserTemplate = "agoltools/templates/unsupportedbrowser.html"
)

type HandlerFunc func(*Request) error

func HandleFunc(pattern string, handler HandlerFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(w, r)

		if err := handler(req); err != nil {
			handleError(req, err)
		}
	})
}

func handleError(r *Request, err error) {
	r.LogInfof("%s", err)

	r.AddData(map[string]interface{}{
		"PageTitle":   "Something went wrong",
		"Error":       err,
		"TryAgainUrl": r.URLWithQuery(),
		"IssuesUrl":   "https://github.com/keyurva/agoltools/issues", // TODO: config?
	})

	if d, ok := err.(DisplayError); ok {
		r.Data["DisplayError"] = d.DisplayError()
	}

	if re, ok := err.(*agolclient.RESTError); ok {
		r.W.WriteHeader(re.Code)
	} else {
		r.W.WriteHeader(http.StatusInternalServerError)
	}

	r.RenderUsingBaseTemplate("agoltools/templates/error.html")
}
