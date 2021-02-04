package storage

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	defaultMode = 0777
)

type fileStorage struct {
	dir string
}

func (f *fileStorage) Read(filename string) (data []byte, err error) {
	return ioutil.ReadFile(path.Join(f.dir, filename))
}

func (f *fileStorage) Write(filename string, data []byte) (err error) {
	err = os.MkdirAll(path.Join(f.dir, filepath.Dir(filename)), defaultMode)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(f.dir, filename), data, defaultMode)
}

func (f *fileStorage) List() (filenameList []string) {
	filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		filenameList = append(filenameList, path)
		return nil
	})
	return filenameList
}

// NewFileStorage :
func NewFileStorage(dir string) Storage {
	return &fileStorage{
		dir: dir,
	}
}
