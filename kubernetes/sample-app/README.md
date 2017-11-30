https://kubernetes.io/docs/tutorials/stateless-application/hello-minikube/

```
eval $(minikube docker-env)
docker build -t hello-node:v1 .
kubectl run hello-node --image=hello-node:v1 --port=8080
kubectl expose deployment hello-node --type=LoadBalancer



docker build -t hello-node:v2 .
kubectl set image deployment/hello-node hello-node=hello-node:v2



kubectl delete service hello-node
kubectl delete deployment hello-node


```