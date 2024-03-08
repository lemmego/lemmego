package fs

import "os"

func CreateFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func CreateFileIfNotExists(filename string) (*os.File, error) {
	if !FileExists(filename) {
		return CreateFile(filename)
	}
	return os.Open(filename)
}

func CreateDir(dirname string) error {
	return os.Mkdir(dirname, 0755)
}

func CreateDirIfNotExists(dirname string) error {
	if !DirExists(dirname) {
		return CreateDir(dirname)
	}
	return nil
}

func DirExists(dirname string) bool {
	stat, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func WriteFileIfNotExists(filename string, data []byte) error {
	if !FileExists(filename) {
		return WriteFile(filename, data)
	}
	return nil
}
