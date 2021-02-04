package storage

import (
	"errors"
)

type memStorage struct {
	data map[string][]byte
}

func (m *memStorage) Read(filename string) (data []byte, err error) {
	data, ok := m.data[filename]
	if !ok {
		return nil, errors.New("File Not Found")
	}
	return data, nil
}

func (m *memStorage) Write(filename string, data []byte) (err error) {
	m.data[filename] = data
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
		data: make(map[string][]byte),
	}
}
