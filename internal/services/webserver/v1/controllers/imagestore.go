package controllers

import (
	"io"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/etag"
)

type ImagestoreController struct {
	st storage.Storage
}

func (c *ImagestoreController) Setup(container di.Container, router fiber.Router) {
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)

	router.Get("/:id", c.getImage)
}

func (c *ImagestoreController) getImage(ctx *fiber.Ctx) error {
	path := ctx.Params("id")

	pathSplit := strings.Split(path, ".")
	imageIDstr := pathSplit[0]

	imageID, err := snowflake.ParseString(imageIDstr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid snowflake ID")
	}

	reader, size, err := c.st.GetObject(static.StorageBucketImages, imageID.String())
	if err != nil {
		return fiber.ErrBadRequest
	}

	defer reader.Close()

	img := new(imgstore.Image)

	img.Size = int(size)
	img.Data = make([]byte, img.Size)
	_, err = reader.Read(img.Data)
	if err != nil && err != io.EOF {
		return err
	}

	img.MimeType = mimetype.Detect(img.Data).String()

	etag := etag.Generate(img.Data, false)

	ctx.Set("Content-Type", img.MimeType)
	// 30 days browser caching
	ctx.Set("Cache-Control", "public, max-age=2592000, immutable")
	ctx.Set("ETag", etag)
	return ctx.Send(img.Data)
}
