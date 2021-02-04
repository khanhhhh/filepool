package crypto

type plainDecryptor struct {
}

func (d *plainDecryptor) Decrypt(dataIn []byte) (dataOut []byte, err error) {
	return dataIn, nil
}

func (d *plainDecryptor) Encrypt(dataIn []byte) (dataOut []byte, err error) {
	return dataIn, nil
}

// NewPlainDecryptor :
func NewPlainDecryptor() Decryptor {
	return &plainDecryptor{}
}
