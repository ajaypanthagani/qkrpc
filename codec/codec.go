package codec

import (
	"github.com/quic-go/quic-go"
)

type Codec interface {
	Write(stream *quic.Stream, payload any) error
	Read(stream *quic.Stream, payload any) error
}
