package storage

import "github.com/minio/minio-go"

type Minio struct {
	client *minio.Client
}
