package main

import (
	"context"
	"fmt"
	pb "github.com/dntam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
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
	fmt.Printf("server %v receive message from client\n", s.serverId)
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
	go serve("50052")
	go serve("50053")
	go serve("50054")
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
