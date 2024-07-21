package fsys

import (
	"io"
	"lemmego/api/logger"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

// LocalStorage is an implementation of StorageInterface for local file system.
type LocalStorage struct {
	// Root directory of the storage.
	RootDirectory string
}

func NewLocalStorage(basePath string) *LocalStorage {
	if basePath == "" {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	return &LocalStorage{
		RootDirectory: basePath,
	}
}

func (ls *LocalStorage) Read(path string) (io.ReadCloser, error) {
	fullPath := ls.RootDirectory + "/" + path
	return os.Open(fullPath)
}

func (ls *LocalStorage) Write(path string, contents []byte) error {
	fullPath := ls.RootDirectory + "/" + path
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(contents)
	return err
}

func (ls *LocalStorage) Delete(path string) error {
	fullPath := ls.RootDirectory + "/" + path
	return os.Remove(fullPath)
}

func (ls *LocalStorage) Exists(path string) (bool, error) {
	fullPath := ls.RootDirectory + "/" + path
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (ls *LocalStorage) Rename(oldPath, newPath string) error {
	oldFullPath := ls.RootDirectory + "/" + oldPath
	newFullPath := ls.RootDirectory + "/" + newPath
	return os.Rename(oldFullPath, newFullPath)
}

func (ls *LocalStorage) Copy(sourcePath, destinationPath string) error {
	sourceFullPath := ls.RootDirectory + "/" + sourcePath
	destinationFullPath := ls.RootDirectory + "/" + destinationPath
	sourceFile, err := os.Open(sourceFullPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationFullPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func (ls *LocalStorage) CreateDirectory(path string) error {
	// For local storage, use os.MkdirAll, which doesn't return an error if the directory already exists.
	fullPath := filepath.Join(ls.RootDirectory, path)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		// If the error indicates that the directory already exists, treat it as success
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (ls *LocalStorage) GetUrl(path string) (string, error) {
	// Construct the URL based on the root directory and the provided path
	fullPath := filepath.Join(ls.RootDirectory, path)
	// Assuming you are serving the files via HTTP
	// return fmt.Sprintf("http://yourdomain.com/%s", fullPath)

	return fullPath, nil
}

func (ls *LocalStorage) Open(path string) (*os.File, error) {
	return os.Open(path)
}

func (ls *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader, dir string) (*os.File, error) {
	if exists, _ := ls.Exists(dir); !exists {
		err := ls.CreateDirectory(dir)
		if err != nil {
			return nil, err
		}
	}

	data, _ := io.ReadAll(file)
	err := ls.Write(path.Join(dir, header.Filename), data)
	if err != nil {
		return nil, err
	}

	if storedFile, err := ls.Open(path.Join(dir, header.Filename)); err != nil {
		return nil, err
	} else {
		defer func() {
			err := storedFile.Close()
			if err != nil {
				logger.V().Info("Uploaded file could not be closed", "Error:", err)
			}
		}()
		return storedFile, nil
	}
}
