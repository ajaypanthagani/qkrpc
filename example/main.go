package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajaypanthagani/qkrpc/example/proto"

	"github.com/ajaypanthagani/qkrpc"

	"github.com/quic-go/quic-go"
)

func main() {
	go runServer()
	time.Sleep(1 * time.Second)
	runClient()
}

func runServer() {
	tlsConfig, err := qkrpc.LoadTLSConfig("keys/cert.pem", "keys/key.pem")
	if err != nil {
		log.Fatal("Failed to load TLS config:", err)
	}

	server := qkrpc.NewQkServer("localhost:4242", tlsConfig)

	// Register an RPC handler
	server.RegisterHandler("echo.EchoService.SayHello", func(ctx context.Context, stream *quic.Stream) error {
		var req proto.HelloRequest
		if err := qkrpc.ReadProtobuf(stream, &req); err != nil {
			return err
		}

		log.Println("Server received:", req.Message)
		resp := &proto.HelloResponse{Reply: "Hello, " + req.Message}
		return qkrpc.WriteProtobuf(stream, resp)
	})

	log.Println("Starting server on :4242")
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

func runClient() {
	tlsConfig, err := qkrpc.LoadClientTLS("keys/cert.pem")
	if err != nil {
		log.Fatal("Failed to load TLS config:", err)
	}

	conn, err := qkrpc.Dial(context.Background(), "localhost:4242", tlsConfig)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	stream, err := conn.Call(context.Background(), "echo.EchoService.SayHello")
	if err != nil {
		log.Fatal("Call failed:", err)
	}

	req := &proto.HelloRequest{Message: "Ajay"}
	if err := qkrpc.WriteProtobuf(stream, req); err != nil {
		log.Fatal("Failed to write request:", err)
	}

	var resp proto.HelloResponse
	if err := qkrpc.ReadProtobuf(stream, &resp); err != nil {
		log.Fatal("Failed to read response:", err)
	}

	fmt.Println("Client received:", resp.Reply)
}
