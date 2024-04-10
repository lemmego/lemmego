package fsys

import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
)

// GCSStorage is an implementation of StorageInterface for Google Cloud Storage.
type GCSStorage struct {
	// GCS bucket name
	BucketName string

	// GCS client
	Client *storage.Client
}

func NewGCSStorage(projectID, bucket string) (*GCSStorage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSStorage{
		BucketName: bucket,
		Client:     client,
	}, nil
}

func (gcs *GCSStorage) Read(path string) (io.ReadCloser, error) {
	ctx := context.Background()
	reader, err := gcs.Client.Bucket(gcs.BucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (gcs *GCSStorage) Write(path string, contents []byte) error {
	ctx := context.Background()
	writer := gcs.Client.Bucket(gcs.BucketName).Object(path).NewWriter(ctx)
	defer writer.Close()

	_, err := writer.Write(contents)
	if err != nil {
		return err
	}
	return nil
}

func (gcs *GCSStorage) Delete(path string) error {
	ctx := context.Background()
	return gcs.Client.Bucket(gcs.BucketName).Object(path).Delete(ctx)
}

func (gcs *GCSStorage) Exists(path string) (bool, error) {
	ctx := context.Background()
	_, err := gcs.Client.Bucket(gcs.BucketName).Object(path).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gcs *GCSStorage) Rename(oldPath, newPath string) error {
	if err := gcs.Copy(oldPath, newPath); err != nil {
		return err
	}
	return gcs.Delete(oldPath)
}

func (gcs *GCSStorage) Copy(sourcePath, destPath string) error {
	ctx := context.Background()
	src := gcs.Client.Bucket(gcs.BucketName).Object(sourcePath)
	dst := gcs.Client.Bucket(gcs.BucketName).Object(destPath)
	_, err := dst.CopierFrom(src).Run(ctx)
	return err
}

func (gcs *GCSStorage) CreateDirectory(path string) error {
	ctx := context.Background()
	obj := gcs.Client.Bucket(gcs.BucketName).Object(path + "/")
	w := obj.NewWriter(ctx)
	if err := w.Close(); err != nil {
		// If the error message indicates that the object already exists, treat it as success
		if strings.Contains(err.Error(), "PreconditionFailed") {
			return nil
		}
		return err
	}
	return nil
}

func (gcs *GCSStorage) GetUrl(path string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcs.BucketName, path)
}
