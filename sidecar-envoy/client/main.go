package main

import (
	"context"
	"fmt"
	pb "github.com/dntam00/grpc-loadbalancing/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	targetConst      = "GRPC_SERVER_ADDR"
	numberOfRequests = 1000
	requestInterval  = 20 * time.Millisecond
)

func main() {
	fmt.Printf("---------------\nstart unary test\n---------------\n")
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go unaryTest(i, &wg)
	}
	wg.Wait()
	fmt.Printf("---------------\nstart stream test\n---------------\n")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go streamTest(i, &wg)
	}
	wg.Wait()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	fmt.Println("finish test client")
}

func unaryTest(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	target := getTarget()
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to server %v fail: %v", target, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewDemoServiceClient(conn)

	responses := make(map[string]int)

	fmt.Printf("start gRPC client %v\n", index)

	for i := 0; i < numberOfRequests; i++ {
		res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "client"})
		if err != nil {
			log.Printf("client %v failed to call SayHello: %v", index, err)
			continue
		}
		responses[res.ServerId] = responses[res.ServerId] + 1
		time.Sleep(requestInterval)
	}
	fmt.Printf("client %v make %v requests, received all response from %v server(s), detail: %+v\n", index, numberOfRequests, len(responses), responses)
}

func getTarget() string {
	target := os.Getenv(targetConst)
	if target == "" {
		panic("target is not set")
	}
	return target
}

func streamTest(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	target := getTarget()
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to server %v fail: %v", target, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewDemoServiceClient(conn)

	responses := make(map[string]int)

	fmt.Printf("start gRPC client %v\n", index)
	stream, err := client.SayHelloStream(context.Background())
	if err != nil {
		log.Fatalf("could not call SayHello: %v", err)
	}

	for i := 0; i < numberOfRequests; i++ {
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
		responses[response.ServerId] = responses[response.ServerId] + 1
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}
	fmt.Printf("client %v make %v requests, received all response from %v server(s), detail: %+v\n", index, numberOfRequests, len(responses), responses)
}
