package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-loadbalancing/model"
	"log"
	"net"
	"os/signal"
	"syscall"
)

type server struct {
	serverId string
	pb.UnimplementedDemoServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("server %v receive message from lb\n", s.serverId)
	return &pb.HelloResponse{ServerId: s.serverId}, nil
}

func main() {
	go serve("50051")
	go serve("50052")
	go serve("50053")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func serve(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDemoServiceServer(s, &server{serverId: port})

	fmt.Println("server is running on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
