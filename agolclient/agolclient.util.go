package agolclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

func randString(n int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	p := make([]rune, n)
	for i := range p {
		p[i] = chars[r.Intn(len(chars))]
	}
	return string(p)
}

func postAndUnmarshalJson(rt http.RoundTripper, url string, params url.Values, v interface{}) (err error) {
	body, err := post(rt, url, params)
	if err != nil {
		return err
	}

	if err = unmarshalJson(body, v); err != nil {
		return err
	}

	return nil
}

func getAndUnmarshalJson(rt http.RoundTripper, url string, params url.Values, v interface{}) (err error) {
	body, err := get(rt, url, params)
	if err != nil {
		return err
	}

	if err = unmarshalJson(body, v); err != nil {
		return err
	}
	return nil
}

func unmarshalJson(body []byte, v interface{}) (err error) {
	errRes := struct {
		Error *RESTError
	}{}

	json.Unmarshal(body, &errRes)

	if errRes.Error != nil {
		return errRes.Error
	}

	return json.Unmarshal(body, v)
}

func get(rt http.RoundTripper, url string, params url.Values) (body []byte, err error) {
	var requestUrl string
	if params != nil {
		requestUrl = fmt.Sprintf("%s?%s", url, params.Encode())
	} else {
		requestUrl = url
	}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	return roundTrip(rt, req)
}

func post(rt http.RoundTripper, url string, params url.Values) (body []byte, err error) {
	var reader io.Reader
	if params != nil {
		reader = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return roundTrip(rt, req)
}

func roundTrip(rt http.RoundTripper, req *http.Request) (body []byte, err error) {
	resp, err := rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
