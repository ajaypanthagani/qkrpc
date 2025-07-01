package compression

import (
	"github.com/golang/snappy"
)

type snappyCompressor struct{}

func NewSnappyCompressor() Compressor {
	return &snappyCompressor{}
}

func (snappyCompressor) Compress(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (snappyCompressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

func (snappyCompressor) Name() string {
	return "snappy"
}
