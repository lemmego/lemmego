package fsys

import (
	"io"
	"mime/multipart"
	"os"
)

// FS defines the methods that any storage system must implement.
type FS interface {
	// Read a file from storage.
	Read(path string) (io.ReadCloser, error)

	// Write a file to storage.
	Write(path string, contents []byte) error

	// Delete a file from storage.
	Delete(path string) error

	// Check if a file exists in storage.
	Exists(path string) (bool, error)

	// Rename a file in storage.
	Rename(oldPath, newPath string) error

	// Copy a file in storage.
	Copy(sourcePath, destinationPath string) error

	// CreateDirectory creates a new directory if doesn't already exist for the given path
	CreateDirectory(path string) error

	// Get the URL for a file in storage (optional).
	// This method may not be applicable to all storage systems.
	// For example, local storage may return a file path, while cloud storage may return a URL.
	GetUrl(path string) (string, error)

	// Open opens a file
	Open(path string) (*os.File, error)

	// Upload uploads a file to the implemented driver
	Upload(file multipart.File, header *multipart.FileHeader, dir string) (*os.File, error)
}
