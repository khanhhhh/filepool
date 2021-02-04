package crypto

// Encryptor :
type Encryptor interface {
	Encrypt(plainText []byte) (cipherText []byte, err error)
}
