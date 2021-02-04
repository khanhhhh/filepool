package pool

import "github.com/khanhhhh/filepool/storage"

type pool struct {
	client storage.Storage
	server storage.Storage
}
