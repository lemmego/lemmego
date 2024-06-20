// The file cache driver implementation for the cache package.
package cache

import ()

type FileStore struct {
	prefix string
}

func NewFileStore(prefix string) *FileStore {
	return &FileStore{
		prefix: prefix,
	}
}

func (f *FileStore) Get(key string) interface{} {
	return nil
}

func (f *FileStore) Many(keys []string) map[string]interface{} {
	return nil
}

func (f *FileStore) Put(key string, value interface{}, seconds int) {
}

func (f *FileStore) PutMany(values map[string]interface{}, seconds int) {
}

func (f *FileStore) Increment(key string, value int) int {
	return 0
}

func (f *FileStore) Decrement(key string, value int) int {
	return 0
}

func (f *FileStore) Forever(key string, value interface{}) {
}

func (f *FileStore) Forget(key string) bool {
	return true
}

func (f *FileStore) Flush() bool {
	return true
}
