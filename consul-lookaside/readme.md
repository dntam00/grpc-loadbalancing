## Consul

```bash
consul agent -dev
consul services register grpc-server-1.json
consul services deregister grpc-server-2.json
```

Log when registering new gRPC server to consul

```bash
GRPC_GO_LOG_SEVERITY_LEVEL=info
GRPC_GO_LOG_VERBOSITY_LEVEL=99
```


```
2025/06/12 19:50:23 INFO: [core] original dial target is: "consul://localhost:8500/kaixin-service?wait=14s"
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel created
2025/06/12 19:50:23 INFO: [core] [Channel #1]parsed dial target is: resolver.Target{URL:url.URL{Scheme:"consul", Opaque:"", User:(*url.Userinfo)(nil), Host:"localhost:8500", Path:"/kaixin-service", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"wait=14s", Fragment:"", RawFragment:""}}
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel authority set to "kaixin-service"
start client 0
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel exiting idle mode
2025/06/12 19:50:23 INFO: [Consul resolver] 1 endpoints fetched in(+wait) 3.972958ms for target={service='kaixin-service' healthy='false' tag=''}
2025/06/12 19:50:23 INFO: [core] [Channel #1]Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:50052",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Metadata": null
    }
  ],
  "Endpoints": [
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:50052",
          "ServerName": "",
          "Attributes": null,
          "BalancerAttributes": null,
          "Metadata": null
        }
      ],
      "Attributes": null
    }
  ],
  "ServiceConfig": null,
  "Attributes": null
} (resolver returned new addresses)
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel switches to new LB policy "round_robin"
2025/06/12 19:50:23 INFO: [roundrobin] [0x1400017f7a0] Created
2025/06/12 19:50:23 INFO: [core] [Channel #1 SubChannel #2]Subchannel created
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel Connectivity change to CONNECTING
2025/06/12 19:50:23 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to CONNECTING
2025/06/12 19:50:23 INFO: [core] [Channel #1 SubChannel #2]Subchannel picks a new address "127.0.0.1:50052" to connect
2025/06/12 19:50:23 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to READY
2025/06/12 19:50:23 INFO: [core] [Channel #1]Channel Connectivity change to READY
2025/06/12 19:50:34 INFO: [Consul resolver] 2 endpoints fetched in(+wait) 10.669448208s for target={service='kaixin-service' healthy='false' tag=''}
2025/06/12 19:50:34 INFO: [core] [Channel #1]Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:50052",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Metadata": null
    },
    {
      "Addr": "127.0.0.1:50053",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Metadata": null
    }
  ],
  "Endpoints": [
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:50052",
          "ServerName": "",
          "Attributes": null,
          "BalancerAttributes": null,
          "Metadata": null
        }
      ],
      "Attributes": null
    },
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:50053",
          "ServerName": "",
          "Attributes": null,
          "BalancerAttributes": null,
          "Metadata": null
        }
      ],
      "Attributes": null
    }
  ],
  "ServiceConfig": null,
  "Attributes": null
} ()
2025/06/12 19:50:34 INFO: [core] [Channel #1 SubChannel #4]Subchannel created
2025/06/12 19:50:34 INFO: [core] [Channel #1 SubChannel #4]Subchannel Connectivity change to CONNECTING
2025/06/12 19:50:34 INFO: [core] [Channel #1 SubChannel #4]Subchannel picks a new address "127.0.0.1:50053" to connect
2025/06/12 19:50:34 INFO: [core] [Channel #1 SubChannel #4]Subchannel Connectivity change to READY
2025/06/12 19:50:48 INFO: [Consul resolver] 2 endpoints fetched in(+wait) 14.66366325s for target={service='kaixin-service' healthy='false' tag=''}
2025/06/12 19:50:48 INFO: [core] [Channel #1]Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:50052",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Metadata": null
    },
    {
      "Addr": "127.0.0.1:50053",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Metadata": null
    }
  ],
  "Endpoints": [
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:50052",
          "ServerName": "",
          "Attributes": null,
          "BalancerAttributes": null,
          "Metadata": null
        }
      ],
      "Attributes": null
    },
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:50053",
          "ServerName": "",
          "Attributes": null,
          "BalancerAttributes": null,
          "Metadata": null
        }
      ],
      "Attributes": null
    }
  ],
  "ServiceConfig": null,
  "Attributes": null
} ()
```