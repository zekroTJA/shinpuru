package imgstore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
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

	img.ID = snowflakenodes.NodeImages.Generate()

	return img, nil
}

func GetImageLink(ident, publicAddr string) string {
	if ident == "" || strings.HasPrefix(ident, "http://") || strings.HasPrefix(ident, "https://") {
		return ident
	}

	return fmt.Sprintf("%s/imagestore/%s.png", publicAddr, ident)
}
