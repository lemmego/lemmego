package fsys

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// MemoryStorage is an implementation of StorageInterface for in-memory storage.
type MemoryStorage struct {
	// Map to store file contents in memory
	data map[string][]byte
	// Mutex to synchronize access to the data map
	mu sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (ms *MemoryStorage) Read(path string) (io.ReadCloser, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if data, ok := ms.data[path]; ok {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return nil, fmt.Errorf("file not found: %s", path)
}

func (ms *MemoryStorage) Write(path string, contents []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.data[path] = contents
	return nil
}

func (ms *MemoryStorage) Delete(path string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.data, path)
	return nil
}

func (ms *MemoryStorage) Exists(path string) (bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	_, ok := ms.data[path]
	return ok, nil
}

func (ms *MemoryStorage) Rename(oldPath, newPath string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if data, ok := ms.data[oldPath]; ok {
		ms.data[newPath] = data
		delete(ms.data, oldPath)
		return nil
	}
	return fmt.Errorf("file not found: %s", oldPath)
}

func (ms *MemoryStorage) Copy(sourcePath, destPath string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if data, ok := ms.data[sourcePath]; ok {
		ms.data[destPath] = data
		return nil
	}
	return fmt.Errorf("file not found: %s", sourcePath)
}

func (ms *MemoryStorage) GetUrl(path string) string {
	// For in-memory storage, we don't have URLs since it's not accessible via HTTP
	return ""
}

func (ms *MemoryStorage) CreateDirectory(path string) error {
	// For in-memory storage, directories are not relevant
	return nil
}
