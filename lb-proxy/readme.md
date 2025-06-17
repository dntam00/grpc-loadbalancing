## LB proxy

```bash
haproxy -f /opt/homebrew/etc/haproxy.cfg
```

## Test result

### L4 proxy

gRPC debug log:

```log
2025/06/17 22:21:34 INFO: [core] original dial target is: "127.0.0.1:8443"
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel created
2025/06/17 22:21:34 INFO: [core] [Channel #1]parsed dial target is: resolver.Target{URL:url.URL{Scheme:"dns", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/127.0.0.1:8443", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""}}
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel authority set to "127.0.0.1:8443"
2025/06/17 22:21:34 INFO: [core] [Channel #1]Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:8443",
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
          "Addr": "127.0.0.1:8443",
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
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel switches to new LB policy "round_robin"
2025/06/17 22:21:34 INFO: [roundrobin] [0x1400017f9b0] Created
2025/06/17 22:21:34 INFO: [pick-first-leaf-lb] [pick-first-leaf-lb 0x1400017af30] Received new config {
  "shuffleAddressList": false
}, resolver state {
  "Addresses": null,
  "Endpoints": [
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:8443",
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
  "Attributes": {
    "\u003c%!p(pickfirstleaf.enableHealthListenerKeyType={})\u003e": "\u003c%!p(bool=true)\u003e"
  }
}
2025/06/17 22:21:34 INFO: [core] [Channel #1 SubChannel #2]Subchannel created
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel Connectivity change to CONNECTING
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel exiting idle mode
2025/06/17 22:21:34 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to CONNECTING
2025/06/17 22:21:34 INFO: [core] [Channel #1 SubChannel #2]Subchannel picks a new address "127.0.0.1:8443" to connect
2025/06/17 22:21:34 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to READY
2025/06/17 22:21:34 INFO: [pick-first-leaf-lb] [pick-first-leaf-lb 0x1400017af30] SubConn 0x1400003c780 reported connectivity state READY. Registering health listener.
2025/06/17 22:21:34 INFO: [core] [Channel #1]Channel Connectivity change to READY
```

gRPC client test unary:

```log
client 0 make 2000 requests, received responses from 1 server(s), detail: map[50054:2000]
finish test client
```

### L7 proxy

gRPC debug log:

```log
2025/06/17 22:29:10 INFO: [core] original dial target is: "127.0.0.1:8443"
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel created
2025/06/17 22:29:10 INFO: [core] [Channel #1]parsed dial target is: resolver.Target{URL:url.URL{Scheme:"dns", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/127.0.0.1:8443", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""}}
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel authority set to "127.0.0.1:8443"
2025/06/17 22:29:10 INFO: [core] [Channel #1]Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:8443",
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
          "Addr": "127.0.0.1:8443",
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
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel switches to new LB policy "round_robin"
2025/06/17 22:29:10 INFO: [roundrobin] [0x1400012b9b0] Created
2025/06/17 22:29:10 INFO: [pick-first-leaf-lb] [pick-first-leaf-lb 0x14000126f30] Received new config {
  "shuffleAddressList": false
}, resolver state {
  "Addresses": null,
  "Endpoints": [
    {
      "Addresses": [
        {
          "Addr": "127.0.0.1:8443",
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
  "Attributes": {
    "\u003c%!p(pickfirstleaf.enableHealthListenerKeyType={})\u003e": "\u003c%!p(bool=true)\u003e"
  }
}
2025/06/17 22:29:10 INFO: [core] [Channel #1 SubChannel #2]Subchannel created
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel Connectivity change to CONNECTING
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel exiting idle mode
2025/06/17 22:29:10 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to CONNECTING
2025/06/17 22:29:10 INFO: [core] [Channel #1 SubChannel #2]Subchannel picks a new address "127.0.0.1:8443" to connect
2025/06/17 22:29:10 INFO: [core] [Channel #1 SubChannel #2]Subchannel Connectivity change to READY
2025/06/17 22:29:10 INFO: [pick-first-leaf-lb] [pick-first-leaf-lb 0x14000126f30] SubConn 0x14000096730 reported connectivity state READY. Registering health listener.
2025/06/17 22:29:10 INFO: [core] [Channel #1]Channel Connectivity change to READY
```

gRPC client test unary:

```log
client 0 make 2000 requests, received responses from 3 server(s), detail: map[50052:667 50053:666 50054:667]
finish test client
```