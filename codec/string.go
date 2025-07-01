package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/quic-go/quic-go"
)

type stringCodec struct {
	compressor compression.Compressor
}

func NewStringCodec(compressor compression.Compressor) Codec {
	return &stringCodec{
		compressor: compressor,
	}
}

// Write writes a length-prefixed UTF-8 string to a QUIC stream.
func (c *stringCodec) Write(stream *quic.Stream, payload any) error {
	s, ok := payload.(string)
	if !ok {
		return fmt.Errorf("Write expects string, got %T", payload)
	}

	data := []byte(s)

	if c.compressor != nil {
		var err error
		data, err = c.compressor.Compress(data)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
	}

	length := uint32(len(data))

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)

	if _, err := stream.Write(lenBuf); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write string data: %w", err)
	}

	return nil
}

// Read reads a length-prefixed UTF-8 string from a QUIC stream.
func (c *stringCodec) Read(stream *quic.Stream, payload any) error {
	strPtr, ok := payload.(*string)
	if !ok {
		return fmt.Errorf("Read expects *string, got %T", payload)
	}

	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(stream, lenBuf); err != nil {
		return fmt.Errorf("failed to read length prefix: %w", err)
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length == 0 {
		*strPtr = ""
		return nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return fmt.Errorf("failed to read string payload: %w", err)
	}

	if c.compressor != nil {
		var err error
		data, err = c.compressor.Decompress(data)
		if err != nil {
			return fmt.Errorf("decompression failed: %w", err)
		}
	}

	*strPtr = string(data)
	return nil
}
