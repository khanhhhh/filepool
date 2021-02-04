package storage

// Storage : storage interface
type Storage interface {
	Read(filename string) (data []byte, err error)
	Write(filename string, data []byte) (err error)
}

// Error :
type Error struct {
}
