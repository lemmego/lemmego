package fsys

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Storage is an implementation of StorageInterface for Amazon S3.
type S3Storage struct {
	// S3 bucket name
	BucketName string

	// AWS session
	Session *session.Session

	// AWS S3 client
	S3Client *s3.S3
}

func NewS3Storage(region string, bucket string) *S3Storage {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		// Add any other necessary configuration options here
	})

	if err != nil {
		panic(err)
	}

	// Create S3 client
	s3Client := s3.New(sess)

	// Initialize S3 storage
	return &S3Storage{
		BucketName: bucket,
		Session:    sess,
		S3Client:   s3Client,
	}
}

func (s3s *S3Storage) Read(path string) (io.ReadCloser, error) {
	// Specify the bucket name and object key
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3s.BucketName),
		Key:    aws.String(path),
	}

	// Retrieve the object
	result, err := s3s.S3Client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// Write writes the given contents to the specified path in S3 storage.
func (s3s *S3Storage) Write(path string, contents []byte) error {
	// Specify the bucket name, object key, and content
	input := &s3.PutObjectInput{
		Bucket: aws.String(s3s.BucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(contents),
	}

	// Upload the object to S3
	_, err := s3s.S3Client.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the file at the specified path from S3 storage.
func (s3s *S3Storage) Delete(path string) error {
	// Specify the bucket name and object key
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3s.BucketName),
		Key:    aws.String(path),
	}

	// Delete the object from S3
	_, err := s3s.S3Client.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if the file exists at the specified path in S3 storage.
func (s3s *S3Storage) Exists(path string) (bool, error) {
	// Specify the bucket name and object key
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s3s.BucketName),
		Key:    aws.String(path),
	}

	// HeadObject checks if the object exists without retrieving the object itself
	_, err := s3s.S3Client.HeadObject(input)
	if err != nil {
		// If the error is NoSuchKey, it means the object doesn't exist
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		// If any other error occurs, return it
		return false, err
	}

	// If no error occurred, the object exists
	return true, nil
}

// Rename renames the file from the oldPath to the newPath in S3 storage.
func (s3s *S3Storage) Rename(oldPath, newPath string) error {
	// Use CopyObject to copy the object to the new path
	_, err := s3s.S3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s3s.BucketName),
		CopySource: aws.String(s3s.BucketName + "/" + oldPath),
		Key:        aws.String(newPath),
	})
	if err != nil {
		return err
	}

	// After copying, delete the original object
	if err := s3s.Delete(oldPath); err != nil {
		return err
	}

	return nil
}

// Copy copies the file from the source path to the destination path in S3 storage.
func (s3s *S3Storage) Copy(sourcePath, destPath string) error {
	_, err := s3s.S3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s3s.BucketName),
		CopySource: aws.String(s3s.BucketName + "/" + sourcePath),
		Key:        aws.String(destPath),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3s *S3Storage) CreateDirectory(path string) error {
	// S3 doesn't have directories in the same way as local filesystems,
	// but you can mimic directory creation by creating an empty object with a trailing slash.
	input := &s3.PutObjectInput{
		Bucket: aws.String(s3s.BucketName),
		Key:    aws.String(path + "/"), // Add trailing slash to simulate directory
		Body:   bytes.NewReader([]byte{}),
	}

	_, err := s3s.S3Client.PutObject(input)
	if err != nil {
		// If the error indicates that the object already exists, treat it as success
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "ObjectAlreadyExists" {
			return nil
		}
		return err
	}
	return nil
}

// GetUrl returns the URL of the file at the specified path in S3 storage.
func (s3s *S3Storage) GetUrl(path string) (string, error) {
	// Format the URL based on the bucket name and object key
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3s.BucketName, path), nil
}

func (s3s *S3Storage) Open(path string) (*os.File, error) {
	panic("not implemented")
}

func (s3s *S3Storage) Upload(file multipart.File, header *multipart.FileHeader, dir string) (*os.File, error) {
	panic("not implemented")
}
