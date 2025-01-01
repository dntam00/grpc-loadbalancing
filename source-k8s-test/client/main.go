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
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	targetAddressConst         = "GRPC_SERVER_ADDR"
	targetPortConst            = "GRPC_SERVER_PORT"
	clientConnectionConst      = "CLIENT_CONNECTION"
	streamerPerConnectionConst = "STREAMER_PER_CONNECTION"
	requestPerClientConst      = "REQUEST_PER_CLIENT"
)

var (
	target                = ""
	targetAddress         = ""
	targetPort            = ""
	clientConnection      = 0
	streamerPerConnection = 20
	requestPerClient      = 1000
)

func getEnv() {
	targetAddress = os.Getenv(targetAddressConst)
	if targetAddress == "" {
		panic(fmt.Sprintf("%v is not set", targetAddressConst))
	}

	targetPort = os.Getenv(targetPortConst)
	if targetPort == "" {
		panic(fmt.Sprintf("%v is not set", targetPortConst))
	}

	target = fmt.Sprintf("%v:%v", targetAddress, targetPort)

	var err error

	clientConnection, err = strconv.Atoi(os.Getenv(clientConnectionConst))
	if err != nil || clientConnection == 0 {
		panic(fmt.Sprintf("%v is not set", clientConnectionConst))
	}

	streamerPerConnection, err = strconv.Atoi(os.Getenv(streamerPerConnectionConst))
	if err != nil || streamerPerConnection == 0 {
		panic(fmt.Sprintf("%v is not set", streamerPerConnectionConst))
	}

	requestPerClient, err = strconv.Atoi(os.Getenv(requestPerClientConst))
	if err != nil || requestPerClient == 0 {
		panic(fmt.Sprintf("%v is not set", requestPerClientConst))
	}
}

func main() {
	getEnv()

	fmt.Printf("---------------\nstart unary test\n---------------\n")
	wg := sync.WaitGroup{}
	for i := 0; i < clientConnection; i++ {
		wg.Add(1)
		go unaryTest(i, &wg)
	}
	wg.Wait()
	fmt.Printf("---------------\nstart stream test\n---------------\n")
	for i := 0; i < clientConnection; i++ {
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
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	if err != nil {
		log.Printf("connect to server %v fail: %v", target, err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewDemoServiceClient(conn)

	responses := make(map[string]int)

	fmt.Printf("start gRPC client %v\n", index)

	for i := 0; i < requestPerClient; i++ {
		res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "client"})
		if err != nil {
			log.Printf("client %v failed to call SayHello: %v", index, err)
			continue
		}
		responses[res.ServerId] = responses[res.ServerId] + 1
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Printf("client %v make %v requests, received all response from %v server(s), detail: %+v\n", index, requestPerClient, len(responses), responses)
}

func streamTest(clientIndex int, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	if err != nil {
		log.Printf("connect to server %v fail: %v\n", target, err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	streamWg := sync.WaitGroup{}

	client := pb.NewDemoServiceClient(conn)

	for i := 0; i < streamerPerConnection; i++ {
		streamWg.Add(1)
		go func(streamIndex int) {
			defer streamWg.Done()
			responses := make(map[string]int)
			fmt.Printf("start stream %v on client %v\n", streamIndex, clientIndex)
			stream, err := client.SayHelloStream(context.Background())

			defer func() {
				if stream != nil {
					_ = stream.CloseSend()
				}
			}()

			if err != nil {
				log.Printf("could not call SayHelloStream on stream %v of client %v, reason: %v\n", streamIndex, clientIndex, err)
				return
			}

			for j := 0; j < requestPerClient; j++ {
				req := &pb.HelloRequest{
					Name: fmt.Sprintf("client %d", j),
				}
				if err := stream.Send(req); err != nil {
					log.Printf("failed to send request: %v\n", err)
				}
				response, err := stream.Recv()
				if err != nil {
					log.Printf("failed to receive response: %v\n", err)
				}
				responses[response.ServerId] = responses[response.ServerId] + 1
			}
			fmt.Printf("client %v make %v requests, received all response from %v server(s), detail: %+v\n", clientIndex, requestPerClient, len(responses), responses)
		}(i)
	}
	streamWg.Wait()
}
