package http

import (
	// "github.com/svenfuchs/travis-go/opts"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var httpClient = http.DefaultClient

var headers = map[string]string{
	"Content-Type":       "application/json",
	"Travis-Api-Version": "3",
}

// NewRequest builds an http request
func NewRequest(method string, path string) *Request {
	r := Request{Method: method, Path: path}
	return &r
}

// Request represents an http request
type Request struct {
	Method string
	Path   string
}

// Response represents an http response
type Response struct {
	Status  int
	Body    []byte
	Headers http.Header
}

var endpoint = "api.travis-ci.org"

// New builds an http client
func New() *HTTP {
	e := endpoint
	// if endp, ok := opts.Get("endpoint"); ok {
	//   e = endp
	// }
	return &HTTP{endpoint: e}
}

// HTTP represents an http client
type HTTP struct {
	endpoint string
}

// Run runs an http request
func (c HTTP) Run(r *Request) (*Response, error) {
	if r.Method == "GET" {
		return c.Get(r.Path, nil)
	}
	return nil, nil
}

// Get makes GET request
func (c HTTP) Get(path string, params *map[string]string) (*Response, error) {
	return c.request("get", path, params)
}

func (c HTTP) request(method string, path string, params *map[string]string) (*Response, error) {
	method = strings.ToUpper(method)
	url := c.urlFor(path, params)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = nil
	if res.StatusCode/100 != 2 {
		err = errors.New(method + " " + url + " " + string(res.Status) + ": " + string(body[:]))
	}

	return &Response{Status: res.StatusCode, Body: body, Headers: res.Header}, err
}

func (c HTTP) urlFor(path string, params *map[string]string) string {
	url := "https://" + c.endpoint + path
	if params != nil {
		path = c.addQuery(path, params)
	}
	return url
}

func (c HTTP) addQuery(path string, params *map[string]string) string {
	uri, err := url.ParseRequestURI(path)
	if err != nil {
		log.Fatal(err)
	}
	q := uri.Query()
	for key, value := range *params {
		q.Set(key, value)
	}
	uri.RawQuery = q.Encode()
	return uri.String()
}
