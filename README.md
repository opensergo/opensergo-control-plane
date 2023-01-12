# OpenSergo Control Plane

OpenSergo control plane enables unified management for microservice governance rules with OpenSergo CRD (under Kubernetes).

![arch](https://user-images.githubusercontent.com/9434884/182856237-8ce85f41-1a1a-4a2a-8f58-db042bd4db42.png)

# How to Start

## step 1 : prepare kubernetes
Please ensure that one of the following conditions is satisfied:
- make sure you are in the kubernetes cluster.
- make sure the `kubeconfig` is ok in you local.

## step 2 : init OpenSergo (CRDs & Namespace & RBAC)
### install all CRDs in following directories:
- [./k8s/crd/bases](./k8s/crd/bases)

you can just only execute commands，like :
``` shell
kubectl apply -f ./k8s/crd/bases/fault-tolerance.opensergo.io_circuitbreakerstrategies.yaml
kubectl apply -f ./k8s/crd/bases/fault-tolerance.opensergo.io_concurrencylimitstrategies.yaml
kubectl apply -f ./k8s/crd/bases/fault-tolerance.opensergo.io_faulttolerancerules.yaml
kubectl apply -f ./k8s/crd/bases/fault-tolerance.opensergo.io_ratelimitstrategies.yaml
kubectl apply -f ./k8s/crd/bases/fault-tolerance.opensergo.io_throttlingstrategies.yaml
kubectl apply -f ./k8s/crd/bases/traffic.opensergo.io_trafficerouters.yaml
```

### install Namespace in following directories:
- [./k8s/namespace.yaml](./k8s/namespace.yaml)

you can just only execute commands，like :
``` shell
kubectl apply -f ./k8s/namespace.yaml
```

### install RBAC in following directories:
- [./k8s/rbac](./k8s/rbac)

you can just only execute commands，like :
``` shell
kubectl apply -f ./k8s/rbac/rbac.yaml
```

### or you can directly execute the file [./cmd/install/init.sh](./cmd/install/init.sh) to install CRDs、 Namespace、 RBAC, like:
``` shell
wget --no-check-certificate https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/cmd/install/init.sh && chmod +x init.sh && ./init.sh
```
if you execute cmd above, some resources will download in `$HOME/opensergo/opensergo-control-plane`.

## step 3 : start OpenSergo Control Plane
we provide several ways to start OpenSergo Control Plane.

The default **port** is **10246**

### 1st :  start it by exec the binary executable file.
download or build the binary executable file
- download the binary executable file in the release or tag page.
- build the binary executable file by yourself.

### 2nd : start it by container image
we provide a container image in Docker Hub, you can start it by the container image.
- Docker Hub: opensergo/opensergo-control-plane:0.0.1-alpha-1
- aliyuncs: opensergo-registry.cn-hangzhou.cr.aliyuncs.com/opensergo/opensergo-control-plane:0.1.0

``` shell
docker run -p 10246:10246 opensergo/opensergo-control-plane
```
If you start it in another container runtime rather than kubernetes pod, you should config the `kubeconfig` in the container runtime.

### 3rd : start it by kubernetes workload (deployment + service)

apply yaml :
- [./cmd/install/k8s/opensergo-control-plane.yaml](./cmd/install/k8s/opensergo-control-plane.yaml)

you can just only execute commands，like :
``` shell
kubectl apply -f ./cmd/install/k8s/opensergo-control-plane.yaml
```

or you can also execute the file [./cmd/install/k8s/deploy.sh](./cmd/install/k8s/deploy.sh) directly, like:
``` shell
wget --no-check-certificate https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/cmd/install/k8s/deploy.sh && chmod +x deploy.sh && ./deploy.sh
```
if you execute cmd above, some resources will download in `$HOME/opensergo/opensergo-control-plane`,   
you can change modify `$HOME/opensergo/opensergo-control-plane/cmd/install/k8s/opensergo-control-plane.yaml` by yourself,   
for example change the Service from `ClusterIP` to `NodePort/LoadBalancer`

### 4th : start it by helm (Work In Process)
we provide a Helm Chart in the repo, you can use it directly.