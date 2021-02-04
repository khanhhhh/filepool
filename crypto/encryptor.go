package crypto

// Encryptor :
type Encryptor interface {
	Encrypt(dataIn []byte) (dataOut []byte, err error)
}
