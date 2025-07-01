package main

import (
	"context"
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
	qkClient := qkrpc.NewQkClient(addr, tlsConfig, protobufCodec)

	err := qkClient.Dial(context.Background())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	req := &proto.HelloRequest{Message: "Ajay"}
	resp := &proto.HelloResponse{}

	err = qkClient.Call(context.Background(), "echo.EchoService.SayHello", req, resp)

	if err != nil {
		log.Fatal("Client received error:", err)
	}

	log.Println("Client received reply: ", resp.Reply)
}
