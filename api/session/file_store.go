package session

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultDir = "storage/session"

type FileStore struct {
	dir string
}

func (fs *FileStore) Delete(token string) error {
	return os.Remove(filepath.Join(fs.dir, token))
}

func (fs *FileStore) Find(token string) ([]byte, bool, error) {
	filename := filepath.Join(fs.dir, token)
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, false, err
	}

	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return nil, false, fmt.Errorf("invalid file format")
	}

	expiry, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, false, err
	}

	if time.Now().After(expiry) {
		os.Remove(filename) // Clean up expired session
		return nil, false, nil
	}

	sessionData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, false, err
	}

	return sessionData, true, nil
}

func (fs *FileStore) Commit(token string, b []byte, expiry time.Time) error {
	data := fmt.Sprintf("%s|%s", expiry.Format(time.RFC3339), base64.StdEncoding.EncodeToString(b))
	return os.WriteFile(filepath.Join(fs.dir, token), []byte(data), 0644)
}

func NewFileStore(directoryPath string) *FileStore {
	if directoryPath == "" {
		directoryPath = defaultDir
	}
	err := os.MkdirAll(directoryPath, 0755)
	if err != nil {
		panic(err)
	}
	return &FileStore{dir: directoryPath}
}
