package codec

import (
	"encoding/binary"
	"io"

	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

// WriteProtobuf sends a length-prefixed protobuf message over a QUIC stream.
func WriteProtobuf(stream *quic.Stream, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
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

// ReadProtobuf reads a length-prefixed protobuf message from a QUIC stream.
func ReadProtobuf(stream *quic.Stream, msg proto.Message) error {
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

	return proto.Unmarshal(data, msg)
}
