package crypto

import "crypto/sha256"

// Constants :
const (
	Length = 32
)

// Hash :
func Hash(data []byte) []byte {
	hashed := sha256.Sum256(data)
	return hashed[:]
}
