package qkrpc

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/ajaypanthagani/qkrpc/codec"
	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/quic-go/quic-go"
)

type QkServer interface {
	Serve() error
	RegisterHandler(name string, handlerFunc func(context.Context, any) any, newReqFunc func() any)
	HandleStream(stream *quic.Stream)
}

func NewQkServer(addr string, tlsConfig *tls.Config, c codec.Codec) QkServer {
	stringCodec := codec.NewStringCodec(compression.NewSnappyCompressor())
	return &qkServer{
		addr:        addr,
		tls:         tlsConfig,
		codec:       c,
		stringCodec: stringCodec,
	}
}

type Handler struct {
	HandlerFunc func(context.Context, any) any
	NewReqFunc  func() any
}

type qkServer struct {
	addr        string
	tls         *tls.Config
	codec       codec.Codec
	stringCodec codec.Codec
	handlers    map[string]Handler
}

// Serve starts the QUIC server and handles incoming connections and streams.
func (s *qkServer) Serve() error {
	listener, err := quic.ListenAddr(s.addr, s.tls, nil)
	if err != nil {
		return err
	}
	log.Printf("qkrpc server listening on %s\n", s.addr)

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go func(conn *quic.Conn) {
			for {
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					log.Println("stream error:", err)
					return
				}
				go s.HandleStream(stream)
			}
		}(conn)
	}
}

// RegisterHandler registers a handler with a name
func (s *qkServer) RegisterHandler(name string, handlerFunc func(context.Context, any) any, newReqFunc func() any) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}

	handler := Handler{
		HandlerFunc: handlerFunc,
		NewReqFunc:  newReqFunc,
	}

	s.handlers[name] = handler
}

// HandleStream reads the method name from the stream and dispatches it to the registered handler.
func (s *qkServer) HandleStream(stream *quic.Stream) {
	var methodName string
	err := s.stringCodec.Read(stream, &methodName)

	if err != nil {
		log.Println("failed to read method name:", err)
		return
	}

	handler, ok := s.handlers[methodName]
	if !ok {
		log.Println("no handler for method:", methodName)
		return
	}

	ctx := context.Background()
	req := handler.NewReqFunc()

	if err := s.codec.Read(stream, req); err != nil {
		log.Println("failed to read req:", err)
		return
	}

	resp := handler.HandlerFunc(ctx, req)

	s.codec.Write(stream, resp)
}
