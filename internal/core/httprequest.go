package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPResponse struct {
	*http.Response
}

func HTTPRequest(method, url string, headers map[string]string, data interface{}) (*HTTPResponse, error) {
	var body io.Reader
	if data != nil {
		var buffer bytes.Buffer
		enc := json.NewEncoder(&buffer)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
		body = bufio.NewReader(&buffer)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return &HTTPResponse{
		Response: resp,
	}, nil
}

func (r *HTTPResponse) BodyAsMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)
	return result, err
}
