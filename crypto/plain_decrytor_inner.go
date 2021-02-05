package crypto

import "io"

type plainDecryptor struct {
}

func (d *plainDecryptor) Decrypt(dataIn []byte) (dataOut []byte, err error) {
	return dataIn, nil
}

func (d *plainDecryptor) DecryptStream(cipherText io.Reader, plainText io.Writer) error {
	_, err := io.Copy(plainText, cipherText)
	return err
}
func (d *plainDecryptor) Encrypt(dataIn []byte) (dataOut []byte, err error) {
	return dataIn, nil
}

func (d *plainDecryptor) EncryptStream(plainText io.Reader, cipherText io.Writer) error {
	_, err := io.Copy(cipherText, plainText)
	return err
}

// NewPlainDecryptor :
func NewPlainDecryptor() Decryptor {
	return &plainDecryptor{}
}
