// Package giphy provides a crappy and inclomplete
// - but at least bloat free - Giphy API client.
package giphy

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

const endpoint = "https://api.giphy.com"

type Client struct {
	apiKey  string
	version string
}

func New(apiKey string, version string) *Client {
	return &Client{apiKey, version}
}

func (c *Client) Search(keyword string, limit, offset int, rating string) (gifs []Gif, err error) {
	var res response[[]Gif]
	err = c.req("GET", "gifs/search", map[string]string{
		"api_key": c.apiKey,
		"q":       keyword,
		"limit":   strconv.Itoa(limit),
		"offset":  strconv.Itoa(offset),
		"rating":  rating,
	}, &res)
	if err != nil {
		return
	}
	gifs = res.Data
	return
}

func (c *Client) req(
	method string,
	path string,
	params map[string]string,
	v any,
) (err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(method)
	req.SetRequestURI(fmt.Sprintf("%s/%s/%s", endpoint, c.version, path))
	for k, v := range params {
		if v != "" {
			req.URI().QueryArgs().Add(k, v)
		}
	}

	if err = fasthttp.Do(req, res); err != nil {
		return
	}

	if res.StatusCode() >= 400 {
		var errModel Error
		json.Unmarshal(res.Body(), &errModel)
		errModel.Code = res.StatusCode()
		err = errModel
		return
	}

	err = json.Unmarshal(res.Body(), v)
	return
}
