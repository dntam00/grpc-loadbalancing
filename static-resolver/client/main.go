package main

import (
	"context"
	"fmt"
	pb "github.com/dntam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/grpclb"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
}

const (
	target = "static:///127.0.0.1:50051,127.0.0.1:50052,127.0.0.1:50053"
)

func main() {
	makeClient()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	fmt.Println("finish test client")
}

func makeClient() {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewDemoServiceClient(conn)

	fmt.Printf("start grpclb test\n")

	for i := 0; i < 10; i++ {
		response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: fmt.Sprintf("client")})
		if err != nil {
			fmt.Println(fmt.Errorf("could not greet: %v", err))
			continue
		}
		fmt.Printf("receive response: %v\n", response.ServerId)
	}
}
