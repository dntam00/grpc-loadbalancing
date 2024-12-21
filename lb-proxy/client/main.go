package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "grpc-loadbalancing/model"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	target          = "127.0.0.1:8443"
	requestInterval = 20 * time.Millisecond
	requestDeadline = 2 * time.Second
)

func main() {
	wg := sync.WaitGroup{}
	stopCh := make(chan struct{})
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go makeClient(i, &wg, stopCh)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	close(stopCh)
	wg.Wait()
	fmt.Println("finish test client")
}

func makeClient(index int, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewDemoServiceClient(conn)

	responses := make(map[string]bool)
	responseLock := sync.Mutex{}

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	requests := 0

	fmt.Printf("start client %v\n", index)

stop:
	for {
		select {
		case <-stopCh:
			break stop
		case <-ticker.C:
			go doSendRequest(index, &requests, c, responses, &responseLock)
		}
	}

	fmt.Printf("client %v make %v requests, received all response from %v server(s), detail: %+v\n", index, requests, len(responses), responses)
}

func doSendRequest(index int, requests *int, c pb.DemoServiceClient, responses map[string]bool, responseLock *sync.Mutex) {
	*requests = *requests + 1
	timeout, cancelFunc := context.WithTimeout(context.Background(), requestDeadline)
	defer cancelFunc()

	// add key for stickiness load balancing at HAProxy
	md := metadata.Pairs("session-id", fmt.Sprintf("unique-session-id-%v", index))
	outgoingContext := metadata.NewOutgoingContext(timeout, md)

	response, err := c.SayHello(outgoingContext, &pb.HelloRequest{Name: fmt.Sprintf("client %v", index)})
	if err != nil {
		fmt.Println(fmt.Errorf("could not greet: %v", err))
		return
	}
	//fmt.Printf("client %v receive response from server %v\n", index, response.ServerId)
	responseLock.Lock()
	defer responseLock.Unlock()
	responses[response.ServerId] = true
}
