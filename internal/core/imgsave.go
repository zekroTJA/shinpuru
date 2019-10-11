package core

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/snowflake"
)

var defClient = http.Client{
	CheckRedirect: func(r *http.Request, via []*http.Request) error {
		r.URL.Opaque = r.URL.Path
		return nil
	},
}

type Image struct {
	ID       snowflake.ID
	MimeType string
	Data     []byte
	Size     int
}

func DownloadImageFromURL(url string) (*Image, error) {
	resp, err := defClient.Get(url)
	if err != nil {
		return nil, err
	}

	img := new(Image)

	img.MimeType = resp.Header.Get("Content-Type")
	if img.MimeType == "" {
		return nil, fmt.Errorf("mime type not received")
	}

	img.Data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img.Size = len(img.Data)

	if img.Data == nil || img.Size < 1 {
		return nil, fmt.Errorf("empty body received")
	}

	return img, nil
}
