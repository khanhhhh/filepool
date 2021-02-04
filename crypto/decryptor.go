package crypto

// Decryptor :
type Decryptor interface {
	Decrypt(dataIn []byte) (dataOut []byte, err error)
	Encrypt(dataIn []byte) (dataOut []byte, err error)
}
