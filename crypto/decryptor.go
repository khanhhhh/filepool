package crypto

// Decryptor :
type Decryptor interface {
	Decrypt(cipherText []byte) (plainText []byte, err error)
	Encrypt(plainText []byte) (cipherText []byte, err error)
}
