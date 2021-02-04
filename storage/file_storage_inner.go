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
	subtractDir := func(path string) (string, error) {
		dirAbsPath, err := filepath.Abs(f.dir)
		if err != nil {
			return "", err
		}
		pathAbsPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return pathAbsPath[len(dirAbsPath)+1:], nil
	}
	filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		filename, _ := subtractDir(path)
		filenameList = append(filenameList, filename)
		return nil
	})
	return filenameList
}

func (f *fileStorage) Delete(filename string) {
	var delFile func(filename string)
	delFile = func(filename string) {
		if filename == f.dir {
			return
		}
		err := os.Remove(filename)
		if err != nil {
			return
		}
		// remove all empty directory
		dir := filepath.Dir(filename)
		dirInfo, err := ioutil.ReadDir(dir)
		if err != nil {
			return
		}
		if len(dirInfo) == 0 {
			delFile(dir)
		}
	}
	delFile(path.Join(f.dir, filename))
}

// NewFileStorage :
func NewFileStorage(dir string) (Storage, error) {
	err := os.MkdirAll(dir, defaultMode)
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		dir: dir,
	}, nil
}
