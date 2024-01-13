package httpreq

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/zekrotja/safepool"
)

var responsePool = safepool.New(func() *Response {
	return &Response{
		Response: fasthttp.AcquireResponse(),
	}
})

// Response extends http.Response with some extra
// utility functions.
type Response struct {
	*fasthttp.Response
}

var _ safepool.ResetState = (*Response)(nil)

// JSON parses the response body data to the
// passed object reference using JSON decoder
// and returns errors occured.
func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body(), v)
}

// Release releases the request instance back to
// the request object pool.
func (r *Response) Release() {
	responsePool.Put(r)
}

func (r *Response) ResetState() {}
