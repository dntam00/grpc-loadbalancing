## gRPC load balancing

### k3d

I use [k3d](https://k3d.io/stable/) to run k3s in local machine.

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