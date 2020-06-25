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

	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, mimeType string) error
	GetObject(bucketName, objectName string) (io.ReadCloser, int64, error)
	DeleteObject(bucketName, objectName string) error
}
