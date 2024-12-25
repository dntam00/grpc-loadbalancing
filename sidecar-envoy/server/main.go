package main

import (
	"context"
	"fmt"
	pb "github.com/dangngoctam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"os"

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

func (s *server) SayHelloStream(stream pb.DemoService_SayHelloStreamServer) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive a request: %v", err)
		}
		fmt.Printf("server %v receive message from lb\n", s.serverId)

		// Send a response back to the client
		res := &pb.HelloResponse{
			ServerId: s.serverId,
		}

		// Send the response to the client
		if err := stream.Send(res); err != nil {
			return fmt.Errorf("failed to send response: %v", err)
		}
	}
}

func main() {
	go serve("50051")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func serve(port string) {
	podName, err := getPodName()

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDemoServiceServer(s, &server{serverId: podName})
	grpc_health_v1.RegisterHealthServer(s, &Health{})

	fmt.Printf("server %v is running on port %v", podName, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getPodName() (string, error) {
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		return "", fmt.Errorf("POD_NAME is not set")
	}
	return podName, nil
}
