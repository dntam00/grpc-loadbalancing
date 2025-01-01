package main

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&StaticBuilder{})
}

type StaticBuilder struct{}

func (sb *StaticBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	endpoints := strings.Split(target.Endpoint(), ",")

	r := &StaticResolver{
		endpoints: endpoints,
		cc:        cc,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (sb *StaticBuilder) Scheme() string {
	return "static"
}
