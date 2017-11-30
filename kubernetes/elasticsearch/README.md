Elastic Search pod configuration

https://github.com/kubernetes/examples/tree/master/staging/elasticsearch

```
kubectl create -f service-account.yaml
kubectl create -f es-svc.yaml
kubectl create -f es-rc.yaml
kubectl create -f rbac.yaml
```