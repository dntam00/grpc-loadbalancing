module consul-lookaside

go 1.23.2

require (
	github.com/dntam00/grpc-loadbalancing/grpc/grpc-code v0.0.0
	github.com/dntam00/grpc-loadbalancing/model v0.0.0-20250612125414-e812d95fc60f
	github.com/mbobakov/grpc-consul-resolver v1.5.3
	google.golang.org/grpc v1.73.0
)

replace github.com/dntam00/grpc-loadbalancing/grpc/grpc-code v0.0.0 => ./../grpc-code

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/hashicorp/consul/api v1.20.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
