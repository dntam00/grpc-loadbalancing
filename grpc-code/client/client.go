package client

import (
	"context"
	"fmt"
	pb "github.com/dntam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
	"time"
)

type TestClient interface {
	TestUnary(requestsPerClient int)
	TestStream(numberOfStreams, requestsPerStream int)
}

type dumpClient struct {
	clients []pb.DemoServiceClient
}

var (
	requestInterval = 10 * time.Millisecond
	requestDeadline = 2 * time.Second
)

func NewGRPCClients(numberOfClient int, resolverScheme, target string, opts ...grpc.DialOption) (TestClient, error) {
	fullAddress := target
	if resolverScheme != "" {
		fullAddress = fmt.Sprintf("%s://%s", resolverScheme, target)
	}

	if len(opts) == 0 {
		opts = append(opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		)
	}

	gRPCClients := make([]pb.DemoServiceClient, 0, numberOfClient)
	for i := 0; i < numberOfClient; i++ {
		conn, err := grpc.NewClient(fullAddress, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC client: %w", err)
		}
		gRPCClients = append(gRPCClients, pb.NewDemoServiceClient(conn))
	}

	client := dumpClient{}
	client.clients = gRPCClients
	return &client, nil
}

func (c *dumpClient) TestUnary(requestsPerClient int) {

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	wg := sync.WaitGroup{}

	for i := 0; i < len(c.clients); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			unaryTest(index, requestsPerClient, c.clients[index])
		}(i)
	}

	wg.Wait()
}

func (c *dumpClient) TestStream(numberOfStreams, requestsPerStream int) {
	if len(c.clients) == 0 {
		log.Fatalf("does not exist client")
	}
	wg := sync.WaitGroup{}
	client := c.clients[0]
	for i := 0; i < numberOfStreams; i++ {
		stream, err := client.SayHelloStream(context.Background())
		if err != nil {
			log.Fatalf("could not call SayHello: %v", err)
		}
		wg.Add(1)
		go func(streamIndex int) {
			defer wg.Done()
			doSendRPCOnStream(streamIndex, requestsPerStream, stream)
		}(i)
	}
	wg.Wait()
}

func doSendRPCOnStream(streamIndex int, requestsPerStream int, stream grpc.BidiStreamingClient[pb.HelloRequest, pb.HelloResponse]) {
	responses := make(map[string]int)
	for i := 0; i < requestsPerStream; i++ {
		req := &pb.HelloRequest{
			Name: fmt.Sprintf("client %d", i),
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send request: %v", err)
		}
		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to receive response: %v", err)
		}
		responses[response.ServerId]++
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}
	fmt.Printf("stream %v make %v requests, received all response from %v server(s), detail: %+v\n", streamIndex, requestsPerStream, len(responses), responses)
}

func unaryTest(index int, requestsPerClient int, c pb.DemoServiceClient) {
	responses := make(map[string]int)
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	failedRequests := 0

	for i := 0; i < requestsPerClient; i++ {
		serverId, err := doSendUnaryRPC(index, c)
		if err != nil {
			failedRequests++
			fmt.Printf("client %v failed to send request: %v\n", index, err)
			continue
		}
		responses[serverId]++
		time.Sleep(requestInterval)
	}
	if failedRequests > 0 {
		fmt.Printf("client %v failed to send %v requests\n", index, failedRequests)
	}
	fmt.Printf("client %v make %v requests, received responses from %v server(s), detail: %+v\n", index, requestsPerClient, len(responses), responses)
}

func doSendUnaryRPC(index int, c pb.DemoServiceClient) (string, error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), requestDeadline)
	defer cancelFunc()

	// add key for stickiness load balancing if needed
	md := metadata.Pairs("session-id", fmt.Sprintf("unique-session-id-%v", index))
	outgoingContext := metadata.NewOutgoingContext(timeout, md)

	response, err := c.SayHello(outgoingContext, &pb.HelloRequest{Name: fmt.Sprintf("client %v", index)})
	if err != nil {
		fmt.Println(fmt.Errorf("could not greet: %v", err))
		return "", err
	}
	return response.ServerId, nil
}
