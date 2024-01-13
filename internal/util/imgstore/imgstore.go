package imgstore

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/pkg/httpreq"
)

// Image wraps metadata and data of an image.
type Image struct {
	ID       snowflake.ID
	MimeType string
	Data     []byte
	Size     int
}

func (img *Image) GenerateID() {
	img.ID = snowflakenodes.NodeImages.Generate()
}

// DownloadFromURL tries to GET an image from the
// passed resource URL, downloading it and returning
// the metadata and data of the image as well as
// occured errors.
func DownloadFromURL(url string) (img *Image, err error) {
	body, contentType, err := httpreq.GetFile(url, nil)
	if err != nil {
		return nil, err
	}

	img = new(Image)

	img.MimeType = contentType
	if img.MimeType == "" {
		return nil, fmt.Errorf("mime type not received")
	}

	img.Data, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	img.Size = len(img.Data)

	if img.Data == nil || img.Size < 1 {
		return nil, fmt.Errorf("empty body received")
	}

	img.GenerateID()

	return img, nil
}

// GetLink returns the publicly accessable link for an image
// resource by passed ident and publicAddr.
func GetLink(ident, publicAddr string) string {
	if ident == "" || strings.HasPrefix(ident, "http://") || strings.HasPrefix(ident, "https://") {
		return ident
	}

	return fmt.Sprintf("%s/imagestore/%s.png", publicAddr, ident)
}
