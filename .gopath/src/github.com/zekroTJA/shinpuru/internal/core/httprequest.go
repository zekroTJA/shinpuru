package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type HTTPResponse struct {
	*http.Response
}

func HTTPRequest(method, url string, headers map[string]string, data interface{}) (*HTTPResponse, error) {
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

	return &HTTPResponse{
		Response: resp,
	}, nil
}

func (r *HTTPResponse) ParseJSONBody(v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	return err
}

func HTTPGetFile(uri string) (io.Reader, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
