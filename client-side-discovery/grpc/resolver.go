package custom_grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/resolver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sync"
)

const (
	namespace = "default"
)

type CustomBuilder struct{}

func (sb *CustomBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	grpcServer := os.Getenv("GRPC_SERVER")
	if grpcServer == "" {
		panic("GRPC_SERVER is not set")
	}

	endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), grpcServer, metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("error getting endpoints: %v", err))
	}

	ipAddress := make([]string, 0)
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			ipAddress = append(ipAddress, address.IP)
		}
	}

	fmt.Printf("ipAddress on build: %+v\n", ipAddress)

	r := &StaticResolver{
		endpoints: ipAddress,
		cc:        cc,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (sb *CustomBuilder) Scheme() string {
	return "custom"
}

type StaticResolver struct {
	endpoints []string
	cc        resolver.ClientConn
	sync.Mutex
}

func (r *StaticResolver) ResolveNow(opts resolver.ResolveNowOptions) {
	r.Lock()
	r.doResolve()
	r.Unlock()
}

func (r *StaticResolver) Close() {
}

func (r *StaticResolver) doResolve() {
	var addrs []resolver.Address
	for i, addr := range r.endpoints {
		addrs = append(addrs, resolver.Address{
			Addr:       addr,
			ServerName: fmt.Sprintf("instance-%d", i+1),
		})
	}

	newState := resolver.State{
		Addresses: addrs,
	}

	err := r.cc.UpdateState(newState)
	if err != nil {
		log.Printf("failed to update state: %v", err)
	}
}
