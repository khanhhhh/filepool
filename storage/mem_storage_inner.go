package storage

import (
	"errors"
	"fmt"
)

type fileObject struct {
	data []byte
}

type memStorage struct {
	data map[string]fileObject
}

func (m *memStorage) Read(filename string) (data []byte, err error) {
	fmt.Printf("read: %s\n", filename)
	o, ok := m.data[filename]
	if !ok {
		return nil, errors.New("File Not Found")
	}
	data = o.data
	return data, nil
}
func (m *memStorage) Write(filename string, data []byte) (err error) {
	fmt.Printf("write: %s\n", filename)
	o := fileObject{
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

func (m *memStorage) Delete(filename string) {
	fmt.Printf("delete: %s\n", filename)
	delete(m.data, filename)
}

// NewMemStorage :
func NewMemStorage() Storage {
	return &memStorage{
		data: make(map[string]fileObject),
	}
}
