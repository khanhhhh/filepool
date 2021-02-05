package util

import "io"

const (
	chunkSize = 256 * 1024 * 1024
)

type transform struct {
	readerIn   io.Reader
	applyChunk func(dataIn []byte) (dataOut []byte, err error)
	buffer     []byte
}

func (t *transform) Read(buffer []byte) (n int, err error) {
	if len(t.buffer) >= len(buffer) {
		// if internal buffer larger than or equal to buffer
		n = len(buffer)
		// copy
		copy(buffer, t.buffer)
		// remove copied data from internal buffer
		t.buffer = t.buffer[len(buffer):]
		return n, nil
	}
	// if internal buffer smaller than buffer
	// read chunkSize bytes from readerIn
	buf := make([]byte, chunkSize)
	nn, err := t.readerIn.Read(buf)
	if err == io.EOF {
		// if reader is empty
		if len(t.buffer) <= 0 {
			// if buffer empty, yield nothing
			return 0, io.EOF
		} else {
			// if buffer not empty, yield last chunk of data
			n = len(t.buffer)
			copy(buffer, t.buffer)
			t.buffer = nil
			return n, nil
		}
	}
	// update buffer
	tbuf, err := t.applyChunk(buf[:nn])
	if err != nil {
		// transform error
		return 0, err
	}
	t.buffer = append(t.buffer, tbuf...)
	return t.Read(buffer)
}

// NewTransform :
func NewTransform(
	readerIn io.Reader,
	applyChunk func(dataIn []byte) (dataOut []byte, err error),
) io.Reader {
	return &transform{
		readerIn:   readerIn,
		applyChunk: applyChunk,
		buffer:     make([]byte, 0),
	}
}
