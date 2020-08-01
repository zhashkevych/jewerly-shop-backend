package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, bucket string, object *StorageObject) error
	Get(ctx context.Context, name, bucket string) (*StorageObject, error)
	Delete(ctx context.Context, name, bucket string) error
	CreateBucket(bucketName, location string) error
}

type StorageFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type StorageObject struct {
	Name       string
	Size       int64
	Reader     io.Reader
	RemoteFile StorageFile
}
