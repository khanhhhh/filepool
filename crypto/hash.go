package crypto

import "crypto/sha256"

// Hasher :
type Hasher interface {
	Hash(dataIn []byte) (dataOut []byte)
}

type defaultHasher struct{}

func (h *defaultHasher) Hash(dataIn []byte) (dataOut []byte) {
	hashed := sha256.Sum256(dataIn)
	return hashed[:]
}

// NewHasher :
func NewHasher() Hasher {
	return &defaultHasher{}
}
