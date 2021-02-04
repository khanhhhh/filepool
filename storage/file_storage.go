package storage

import (
	"io/ioutil"
	"path"
)

type fileStorage struct {
	dir string
}

func (f *fileStorage) Read(filename string) (data []byte, err error) {
	return ioutil.ReadFile(path.Join(f.dir, filename))
}

func (f *fileStorage) Write(filename string, data []byte) (err error) {
	return ioutil.WriteFile(path.Join(f.dir, filename), data, 0644)
}

// NewFileStorage :
func NewFileStorage(dir string) Storage {
	return &fileStorage{
		dir: dir,
	}
}
