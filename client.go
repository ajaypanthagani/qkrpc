package qkrpc

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ajaypanthagani/qkrpc/codec"
	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/quic-go/quic-go"
)

type QkClient interface {
	Dial(ctx context.Context) error
	Call(ctx context.Context, method string) (*quic.Stream, error)
}

func NewQkClient(addr string, tlsConfig *tls.Config) QkClient {
	stringCodec := codec.NewStringCodec(compression.NewSnappyCompressor())
	return &qkClient{
		addr:        addr,
		tslConfig:   tlsConfig,
		stringCodec: stringCodec,
	}
}

type qkClient struct {
	addr        string
	tslConfig   *tls.Config
	stringCodec codec.Codec
	conn        *quic.Conn
}

// Dial establishes a QUIC connection to the given address using TLS config.
func (c *qkClient) Dial(ctx context.Context) error {
	conn, err := quic.DialAddr(ctx, c.addr, c.tslConfig, nil)

	if err != nil {
		return fmt.Errorf("couldn't establish a connection to server: %w", err)
	}

	c.conn = conn

	return nil
}

// Call opens a QUIC stream and writes the RPC method name as a length-prefixed string.
func (c *qkClient) Call(ctx context.Context, method string) (*quic.Stream, error) {
	stream, err := c.conn.OpenStreamSync(ctx)

	if err != nil {
		return nil, err
	}

	if err := c.stringCodec.Write(stream, method); err != nil {
		return nil, err
	}

	return stream, nil
}
