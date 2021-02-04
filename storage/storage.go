package storage

import "time"

// Info :
type Info struct {
	Name    string    `json:"name"`
	ModTime time.Time `json:"modtime"`
}

// Storage : storage interface
type Storage interface {
	Read(filename string) (data []byte, err error)
	Stat(filename string) (info Info, err error)
	Write(filename string, data []byte) (err error)
	List() (filenameList []string)
}

// Error :
type Error struct {
}
