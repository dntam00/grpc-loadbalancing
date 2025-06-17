package main

import (
	"fmt"
	gc "github.com/dntam00/grpc-loadbalancing/grpc/grpc-code/client"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"log"
)

const (
	scheme = "consul"
)

func main() {
	client, err := gc.NewGRPCClients(1, scheme, "127.0.0.1:8500/kaixin-service?wait=14s")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}

	//client.TestUnary(2000)

	client.TestStream(100, 1000)

	fmt.Println("finish test client")
}
