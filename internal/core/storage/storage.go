package storage

import (
	"io"

	"github.com/zekroTJA/shinpuru/internal/core/config"
)

type Storage interface {
	Connect(cfg *config.Config) error

	BucketExists(name string) (bool, error)
	CreateBucket(name string, location ...string) error
	CreateBucketIfNotExists(name string, location ...string) error

	PutObject(bucketName string, objectName string, reader io.Reader, objectSize int64, mimeType string) error
	GetObject(bucketName string, objectName string) (io.ReadCloser, error)
}
