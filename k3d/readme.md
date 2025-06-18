## k3d

Example task command to start/delete k8s cluster.

```bash
task create
task delete
```

Update host file to resolve registry domain when working with docker image/

```bash
vim /etc/hosts

127.0.0.1 kaixin-registry
```

Command to get repository metadata.

```bash
curl kaixin-registry:5001/v2/_catalog
curl kaixin-registry:5001/v2/<repo_name>/tags/list
```

## Istio

Setup istio in k8s cluster created by k3d.

```bash
kubectl create namespace istio-system
helm install istio-base istio/base -n istio-system --wait
helm install istiod istio/istiod -n istio-system --wait
kubectl label namespace default istio-injection=enabled --overwrite

# ingress gateway
helm install istio-ingressgateway istio/gateway -n istio-system --wait
```