package storage

import "io"

// Storage : storage interface
type Storage interface {
	Read(filename string) (reader io.ReadCloser, err error)
	Write(filename string) (writer io.WriteCloser, err error)
	List() (filenameList []string)
	Delete(filename string)
}

// Error :
type Error struct {
}
