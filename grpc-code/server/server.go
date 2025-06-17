package server

import (
	"context"
	"fmt"
	pb "github.com/dntam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct {
	serverId string
	pb.UnimplementedDemoServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("server %v receive message from client\n", s.serverId)
	time.Sleep(5 * time.Millisecond)
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

		time.Sleep(5 * time.Millisecond)

		// Send the response to the client
		if err := stream.Send(res); err != nil {
			return fmt.Errorf("failed to send response: %v", err)
		}
	}
}

func Serve(port string) {
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
