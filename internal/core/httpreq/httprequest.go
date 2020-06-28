package httpreq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Response struct {
	*http.Response
}

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

func Get(url string, headers map[string]string) (*Response, error) {
	return Request("GET", url, headers, nil)
}

func Post(url string, headers map[string]string, data interface{}) (*Response, error) {
	return Request("POST", url, headers, data)
}

func (r *Response) JSON(v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	return err
}

func GetFile(uri string) (io.Reader, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
