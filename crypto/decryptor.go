package crypto

import (
	"io"
)

// Decryptor :
type Decryptor interface {
	Decrypt(cipherText []byte) (plainText []byte, err error)
	Encrypt(plainText []byte) (cipherText []byte, err error)
}

const (
	chunk = 256 * 1024 * 1024
)

// TransformStream :
func TransformStream(
	transform func(inText []byte) (outText []byte, err error),
	reader io.Reader,
	writer io.Writer,
) error {
	buffer := make([]byte, chunk)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		writeText, err := transform(buffer[:n])
		if err != nil {
			return err
		}
		n, err = writer.Write(writeText)
		if err != nil {
			return err
		}
	}
	return nil
}
