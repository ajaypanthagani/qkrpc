package qkrpc

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

type ClientConn struct {
	conn *quic.Conn
}

// Dial establishes a QUIC connection to the given address using TLS config.
func Dial(ctx context.Context, addr string, tlsConfig *tls.Config) (*ClientConn, error) {
	conn, err := quic.DialAddr(ctx, addr, tlsConfig, nil)

	if err != nil {
		return nil, err
	}

	return &ClientConn{conn: conn}, nil
}

// OpenStream opens a new bidirectional QUIC stream to send/receive RPC data.
func (c *ClientConn) OpenStream(ctx context.Context) (*quic.Stream, error) {
	return c.conn.OpenStreamSync(ctx)
}

// Call opens a QUIC stream and writes the RPC method name as a length-prefixed string.
func (c *ClientConn) Call(ctx context.Context, method string) (*quic.Stream, error) {
	stream, err := c.OpenStream(ctx)

	if err != nil {
		return nil, err
	}

	if err := WriteString(stream, method); err != nil {
		return nil, err
	}

	return stream, nil
}
