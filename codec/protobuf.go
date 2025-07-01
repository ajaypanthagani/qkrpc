package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

type protobufCodec struct {
	compressor compression.Compressor
}

func NewProtobufCodec(compressor compression.Compressor) Codec {
	return &protobufCodec{
		compressor: compressor,
	}
}

// Write sends a length-prefixed protobuf message over a QUIC stream.
func (c *protobufCodec) Write(stream *quic.Stream, payload any) error {
	msg, ok := payload.(proto.Message)
	if !ok {
		return fmt.Errorf("Write expects proto.Message, got %T", payload)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if c.compressor != nil {
		data, err = c.compressor.Compress(data)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
	}

	// First write the length (4 bytes, big endian)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))

	if _, err := stream.Write(length); err != nil {
		return err
	}

	_, err = stream.Write(data)
	return err
}

// Read reads a length-prefixed protobuf message from a QUIC stream.
func (c *protobufCodec) Read(stream *quic.Stream, payload any) error {
	msg, ok := payload.(proto.Message)
	if !ok {
		return fmt.Errorf("Read expects proto.Message, got %T", payload)
	}

	// Read the first 4 bytes for the message length
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(stream, lengthBuf); err != nil {
		return err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	data := make([]byte, length)

	if _, err := io.ReadFull(stream, data); err != nil {
		return err
	}

	if c.compressor != nil {
		var err error
		data, err = c.compressor.Decompress(data)
		if err != nil {
			return fmt.Errorf("decompression failed: %w", err)
		}
	}

	return proto.Unmarshal(data, msg)
}
