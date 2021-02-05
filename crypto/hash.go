package crypto

import "crypto/sha256"

// Hasher :
type Hasher interface {
	Hash(dataIn []byte) (dataOut []byte, err error)
}

type defaultHasher struct{}

func (h *defaultHasher) Hash(dataIn []byte) (dataOut []byte, err error) {
	hashed := sha256.Sum256(dataIn)
	return hashed[:], nil
}

// NewHasher :
func NewHasher() Hasher {
	return &defaultHasher{}
}
