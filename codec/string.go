package codec

import (
	"encoding/binary"
	"io"

	"github.com/quic-go/quic-go"
)

// WriteString writes a length-prefixed UTF-8 string to a QUIC stream.
func WriteString(stream *quic.Stream, s string) error {
	length := uint32(len(s))
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)

	if _, err := stream.Write(lenBuf); err != nil {
		return err
	}

	_, err := stream.Write([]byte(s))
	return err
}

// ReadString reads a length-prefixed UTF-8 string from a QUIC stream.
func ReadString(stream *quic.Stream) (string, error) {
	lenBuf := make([]byte, 4)

	if _, err := io.ReadFull(stream, lenBuf); err != nil {
		return "", err
	}

	length := binary.BigEndian.Uint32(lenBuf)
	data := make([]byte, length)

	if _, err := io.ReadFull(stream, data); err != nil {
		return "", err
	}

	return string(data), nil
}
