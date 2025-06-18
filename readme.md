## gRPC load balancing

This repository includes demo source code and testing for series gRPC load balancing in this blog: https://notes-ngtam.pages.dev/posts/grpc-load-balancer
- HAProxy as load balancer.
- Service mesh in K8s with Istio and Envoy.
- Lookaside load balancing with xDS Envoy.
- Lookaside load balancing with Consul.
- Headless service in K8s.

### k3d

I use [k3d](https://k3d.io/stable/) to run k3s in local machine.

Please refer to readme file in folder [`k3d`](./k3d) for details setup.

After running, verify cluster with command `k3d cluster list`

```
NAME          SERVERS   AGENTS   LOADBALANCER
local         1/1       2/2      true
```
If you have multiple cluster, you could switch between them with command `kubectl config use-context k3d-<cluster_name>`, and verify by `kubectl config current-context`.

k3d run a container as docker registry, so please make sure the port mapping to host machine is correct.

### go task

I use [go task](https://taskfile.dev/#/) to group multiple commands in one alias.

### proto buf

Run command to generate gRPC code
```bash
    protoc --go_out=. --go-grpc_out=. ./model/service.proto
```
