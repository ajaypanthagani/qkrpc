package compression

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
	Name() string
}
