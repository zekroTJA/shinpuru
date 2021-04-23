package storage

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/zekroTJA/shinpuru/internal/config"
)

// File implements the Storage interface for a
// local file storage provider.
type File struct {
	location string
}

func (f *File) Connect(cfg *config.Config) (err error) {
	f.location = cfg.Storage.File.Location
	return nil
}

func (f *File) BucketExists(name string) (bool, error) {
	stat, err := os.Stat(path.Join(f.location, name))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, errors.New("location is a file")
	}
	return true, nil
}

func (f *File) CreateBucket(name string, location ...string) error {
	return os.MkdirAll(path.Join(f.location, name), os.ModeDir)
}

func (f *File) CreateBucketIfNotExists(name string, location ...string) (err error) {
	ok, err := f.BucketExists(name)
	if err == nil && !ok {
		err = f.CreateBucket(name, location...)
	}

	return
}

func (f *File) PutObject(bucketName string, objectName string, reader io.Reader, objectSize int64, mimeType string) (err error) {
	if err = f.CreateBucketIfNotExists(bucketName); err != nil {
		return
	}

	fd := path.Join(f.location, bucketName, objectName)

	stat, err := os.Stat(fd)
	var fh *os.File

	if os.IsNotExist(err) {
		fh, err = os.Create(fd)
	} else if err != nil {
		return
	} else if stat.IsDir() {
		return errors.New("given file dir is a location")
	} else {
		fh, err = os.Open(fd)
	}

	if err != nil {
		return
	}

	defer fh.Close()

	_, err = io.CopyN(fh, reader, objectSize)
	return
}

func (f *File) GetObject(bucketName string, objectName string) (io.ReadCloser, int64, error) {
	fd := path.Join(f.location, bucketName, objectName)
	stat, err := os.Stat(fd)
	var fh *os.File

	if os.IsNotExist(err) {
		return nil, 0, errors.New("file does not exist")
	} else if err != nil {
		return nil, 0, err
	} else if stat.IsDir() {
		return nil, 0, errors.New("given file dir is a location")
	} else {
		fh, err = os.Open(fd)
	}

	return fh, stat.Size(), err
}

func (f *File) DeleteObject(bucketName, objectName string) error {
	fd := path.Join(f.location, bucketName, objectName)
	return os.Remove(fd)
}
