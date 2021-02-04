package storage

import "errors"

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

// NewMemStorage :
func NewMemStorage() Storage {
	return &memStorage{
		data: make(map[string][]byte),
	}
}
