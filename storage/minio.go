package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	timeout = time.Second*5
)

type FileStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

func NewFileStorage(client *minio.Client, bucket, endpoint string) *FileStorage {
	return &FileStorage{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

// todo: image compression
func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	ctx, clFn := context.WithTimeout(ctx, timeout)
	defer clFn()

	_, err := fs.client.PutObjectWithContext(ctx,
		fs.bucket, input.Name, input.File, input.Size, opts)
	if err != nil {
		logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
		return "", err
	}

	return fs.generateFileURL(input.Name), nil
}

func (fs *FileStorage) generateFileURL(fileName string) string {
	endpoint := strings.Replace(fs.endpoint, "localstack", "localhost", -1)
	return fmt.Sprintf("http://%s/%s/%s", endpoint, fs.bucket, fileName)
}
