package main

import (
	"context"
	"log"
	"time"

	"github.com/ajaypanthagani/qkrpc/codec"
	"github.com/ajaypanthagani/qkrpc/compression"
	"github.com/ajaypanthagani/qkrpc/example/proto"

	"github.com/ajaypanthagani/qkrpc"
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

func newHelloRequestFunc() any {
	return &proto.HelloRequest{}
}

func helloRequestHandlerFunc(ctx context.Context, req any) (resp any) {
	helloReq := req.(*proto.HelloRequest)

	log.Println("Server received:", helloReq.Message)
	return &proto.HelloResponse{Reply: "Hello, " + helloReq.Message}
}

func runServer() {
	tlsConfig, err := qkrpc.LoadTLSConfig("keys/cert.pem", "keys/key.pem")
	if err != nil {
		log.Fatal("Failed to load TLS config:", err)
	}

	server := qkrpc.NewQkServer(addr, tlsConfig, protobufCodec)

	// Register an RPC handler
	server.RegisterHandler("echo.EchoService.SayHello", helloRequestHandlerFunc, newHelloRequestFunc)

	log.Println("Starting server on :4242")
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

func runClient() {
	qkClient := qkrpc.NewQkClient(addr, tlsConfig, protobufCodec)

	err := qkClient.Connect(context.Background())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	req := &proto.HelloRequest{Message: "Ajay"}
	resp := &proto.HelloResponse{}

	// call an RPC endpoint
	err = qkClient.Call(context.Background(), "echo.EchoService.SayHello", req, resp)

	if err != nil {
		log.Fatal("Client received error:", err)
	}

	log.Println("Client received reply: ", resp.Reply)
}
