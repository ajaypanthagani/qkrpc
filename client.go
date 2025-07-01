package qkrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ajaypanthagani/qkrpc/codec"
	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/quic-go/quic-go"
)

type QkClient interface {
	Connect(ctx context.Context) error
	Call(ctx context.Context, method string, request any, response any) error
}

func NewQkClient(addr string, tlsConfig *tls.Config, c codec.Codec) QkClient {
	stringCodec := codec.NewStringCodec(compression.NewSnappyCompressor())
	return &qkClient{
		addr:        addr,
		tslConfig:   tlsConfig,
		stringCodec: stringCodec,
		codec:       c,
	}
}

type qkClient struct {
	addr        string
	tslConfig   *tls.Config
	stringCodec codec.Codec
	codec       codec.Codec
	conn        *quic.Conn
}

// Dial establishes a QUIC connection to the given address using TLS config.
func (c *qkClient) Connect(ctx context.Context) error {
	conn, err := quic.DialAddr(ctx, c.addr, c.tslConfig, nil)

	if err != nil {
		return fmt.Errorf("couldn't establish a connection to server: %w", err)
	}

	c.conn = conn

	return nil
}

// Call opens a QUIC stream and call the RPC method.
func (c *qkClient) Call(ctx context.Context, method string, request any, response any) error {
	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}

	stream, err := c.conn.OpenStreamSync(ctx)

	if err != nil {
		log.Println("Error opening stream to write request:", err)
		return err
	}

	if err := c.stringCodec.Write(stream, method); err != nil {
		log.Println("Error opening stream to write request:", err)
		return err
	}

	if err := c.codec.Write(stream, request); err != nil {
		log.Println("Failed to write request:", err)
		return err
	}

	if err := c.codec.Read(stream, response); err != nil {
		log.Println("Failed to read response:", err)
		return err
	}

	return nil
}
