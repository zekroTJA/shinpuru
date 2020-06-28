// Package httpreq provides general utilities for
// around net/http requests for a simpler API and
// extra utilities for parsing JSON request and
// response boddies.
package httpreq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// Response extends http.Response with some extra
// utility functions.
type Response struct {
	*http.Response
}

// Request executes a HTTP request with the given method to the
// given URL and attaches the passed headers. When data is passed,
// the object will be serialized using JSON encoder and attached to
// the request body.
func Request(method, url string, headers map[string]string, data interface{}) (*Response, error) {
	var body io.Reader
	var dataLen int
	if data != nil {
		var buffer bytes.Buffer
		enc := json.NewEncoder(&buffer)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
		dataLen = buffer.Len()
		body = bufio.NewReader(&buffer)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if dataLen > 0 {
		req.Header.Add("Content-Length", strconv.Itoa(dataLen))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: resp,
	}, nil
}

// Get is shorthand for Request using the GET method.
func Get(url string, headers map[string]string) (*Response, error) {
	return Request("GET", url, headers, nil)
}

// Get is shorthand for Request using the POST method.
func Post(url string, headers map[string]string, data interface{}) (*Response, error) {
	return Request("POST", url, headers, data)
}

// JSON parses the response body data to the
// passed object reference using JSON decoder
// and returns errors occured.
func (r *Response) JSON(v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	return err
}

// GetFile is shorthand for http.Get and returns
// the body as io.Reader as well as occured errors
// during request execution.
func GetFile(uri string) (io.Reader, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
