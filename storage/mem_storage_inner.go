package storage

import (
	"errors"
	"time"
)

type fileObject struct {
	info Info
	data []byte
}

type memStorage struct {
	data map[string]fileObject
}

func (m *memStorage) Read(filename string) (data []byte, err error) {
	o, ok := m.data[filename]
	if !ok {
		return nil, errors.New("File Not Found")
	}
	data = o.data
	return data, nil
}
func (m *memStorage) Stat(filename string) (info Info, err error) {
	o, ok := m.data[filename]
	if !ok {
		return Info{}, errors.New("File Not Found")
	}
	info = o.info
	return info, nil
}

func (m *memStorage) Write(filename string, data []byte) (err error) {
	o := fileObject{
		info: Info{Name: filename, ModTime: time.Now()},
		data: data,
	}
	m.data[filename] = o
	return nil
}

func (m *memStorage) List() (filenameList []string) {
	filenameList = make([]string, 0)
	for filename := range m.data {
		filenameList = append(filenameList, filename)
	}
	return filenameList
}

// NewMemStorage :
func NewMemStorage() Storage {
	return &memStorage{
		data: make(map[string]fileObject),
	}
}
