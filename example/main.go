package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajaypanthagani/qkrpc/codec"
	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/ajaypanthagani/qkrpc/example/proto"

	"github.com/ajaypanthagani/qkrpc"

	"github.com/quic-go/quic-go"
)

func main() {
	go runServer()
	time.Sleep(1 * time.Second)
	runClient()
}

var (
	addr = "localhost:4242"
)

var (
	compressor    = compression.NewSnappyCompressor()
	protobufCodec = codec.NewProtobufCodec(compressor)
	tlsConfig, _  = qkrpc.LoadClientTLS("keys/cert.pem")
)

func runServer() {
	tlsConfig, err := qkrpc.LoadTLSConfig("keys/cert.pem", "keys/key.pem")
	if err != nil {
		log.Fatal("Failed to load TLS config:", err)
	}

	server := qkrpc.NewQkServer(addr, tlsConfig, protobufCodec)

	// Register an RPC handler
	server.RegisterHandler("echo.EchoService.SayHello", func(ctx context.Context, stream *quic.Stream) error {
		var req proto.HelloRequest
		if err := protobufCodec.Read(stream, &req); err != nil {
			return err
		}

		log.Println("Server received:", req.Message)
		resp := &proto.HelloResponse{Reply: "Hello, " + req.Message}
		return protobufCodec.Write(stream, resp)
	})

	log.Println("Starting server on :4242")
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

func runClient() {
	qkClient := qkrpc.NewQkClient(addr, tlsConfig)
	err := qkClient.Dial(context.Background())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	stream, err := qkClient.Call(context.Background(), "echo.EchoService.SayHello")
	if err != nil {
		log.Fatal("Call failed:", err)
	}

	req := &proto.HelloRequest{Message: "Ajay"}

	if err := protobufCodec.Write(stream, req); err != nil {
		log.Fatal("Failed to write request:", err)
	}

	var resp proto.HelloResponse
	if err := protobufCodec.Read(stream, &resp); err != nil {
		log.Fatal("Failed to read response:", err)
	}

	fmt.Println("Client received:", resp.Reply)
}
