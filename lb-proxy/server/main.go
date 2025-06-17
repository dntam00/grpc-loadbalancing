package main

import (
	"context"
	gc "github.com/dntam00/grpc-loadbalancing/grpc/grpc-code/server"
	"os/signal"
	"syscall"
)

func main() {
	go gc.Serve("50052")
	go gc.Serve("50053")
	go gc.Serve("50054")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}
