package qkrpc

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/quic-go/quic-go"
)

type QkServer interface {
	Serve() error
	RegisterHandler(name string, handler func(context.Context, *quic.Stream) error)
	HandleStream(stream *quic.Stream)
}

func NewQkServer(addr string, tlsConfig *tls.Config) QkServer {
	return &qkServer{
		addr: addr,
		tls:  tlsConfig,
	}
}

type qkServer struct {
	addr     string
	tls      *tls.Config
	handlers map[string]func(context.Context, *quic.Stream) error
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
func (s *qkServer) RegisterHandler(name string, handler func(context.Context, *quic.Stream) error) {
	if s.handlers == nil {
		s.handlers = make(map[string]func(context.Context, *quic.Stream) error)
	}
	s.handlers[name] = handler
}

// HandleStream reads the method name from the stream and dispatches it to the registered handler.
func (s *qkServer) HandleStream(stream *quic.Stream) {
	methodName, err := ReadString(stream)

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
	if err := handler(ctx, stream); err != nil {
		log.Println("handler error:", err)
	}
}
