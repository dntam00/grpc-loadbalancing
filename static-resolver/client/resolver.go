package main

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"sync"
)

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

	_ = r.cc.UpdateState(newState)
}
