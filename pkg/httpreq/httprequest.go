// Package httpreq provides general utilities for
// around net/http requests for a simpler API and
// extra utilities for parsing JSON request and
// response boddies.
package httpreq

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/valyala/fasthttp"
)

// Request executes a HTTP request with the given method to the
// given URL and attaches the passed headers. When data is passed,
// the object will be serialized using JSON encoder and attached to
// the request body.
func Request(method, url string, headers map[string]string, data interface{}) (res *Response, err error) {
	defer func() {
		if err != nil && res != nil {
			res.Release()
		}
	}()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res = responsePool.Get().(*Response)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if data != nil {
		err = json.NewEncoder(req.BodyWriter()).Encode(data)
		if err != nil {
			return
		}
	}

	err = fasthttp.Do(req, res.Response)
	return
}

// Get is shorthand for Request using the GET method.
func Get(url string, headers map[string]string) (*Response, error) {
	return Request("GET", url, headers, nil)
}

// Post is shorthand for Request using the POST method.
func Post(url string, headers map[string]string, data interface{}) (*Response, error) {
	return Request("POST", url, headers, data)
}

// GetFile is shorthand for http.Get and returns
// the body as io.Reader as well as occured errors
// during request execution.
func GetFile(url string, headers map[string]string) (r io.Reader, contentType string, err error) {
	resp, err := Get(url, headers)
	if err != nil {
		return
	}
	defer resp.Release()
	r = bytes.NewBuffer(resp.Body())
	contentType = string(resp.Header.Peek("content-type"))
	return
}
