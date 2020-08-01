package minio

import (
	"context"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
	"github.com/zhashkevych/jewelry-shop-backend/storage"
)

// Storage implements storage.Files wraps database connection
type Storage struct {
	client *minio.Client
}

// NewStorage creates new storage
func NewStorage(client *minio.Client) *Storage {
	return &Storage{
		client: client,
	}
}

// Save stores file to storage
func (s *Storage) Save(ctx context.Context, bucket string, object *storage.StorageObject) error {
	n, err := s.client.PutObjectWithContext(ctx, bucket, object.Name,
		object.Reader, object.Size, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "storage")
	}

	if n == 0 {
		return errors.New("storage: write error")
	}

	return nil
}

// Delete removes file from storage
func (s *Storage) Delete(_ context.Context, name, bucket string) error {
	return s.client.RemoveObject(bucket, name)
}

// Get fetches file from storage
func (s *Storage) Get(ctx context.Context, name, bucket string) (*storage.StorageObject, error) {
	obj, err := s.client.GetObjectWithContext(ctx, bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	info, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size == 0 {
		return nil, errors.New("empty file")
	}

	return &storage.StorageObject{
		Name:       name,
		Size:       info.Size,
		RemoteFile: obj,
	}, nil
}

// may be deleted later
func (s *Storage) CreateBucket(bucketName, location string) error {
	err := s.client.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.client.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
