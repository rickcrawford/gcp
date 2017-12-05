esautocomplete
--------------


```bash
docker build -t my-golang-app .
docker run -it --rm --name my-running-app my-golang-app
```


curl http://127.0.0.1:9200/_cat/health -u elastic:changeme




Systems involved

* PaaS
- Google Cloud Storage - store the catalog
- Google Cloud Functions - create a message when updated, record events
- Google Pub Sub - Topics for new-catalog, event, update-records
- App Engine - storage of records, predictor
- IAM - permissions
- Big Query - analysis of recommendations

* IaaS
- Cloud KMS - key management for JWT
- Kubernetes engine - container managemnet
- Stackdriver
- Container Registry
- Cloud CDN
- Cloud SQL
- Cloud Persistent Disk (possible)



TODO:
Use Kubernetes config map to manage secrets and environment variables for elastic/redis


kubectl create secret generic pubsub-key --from-file=key.json=<PATH-TO-KEY-FILE>.json
https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform

# certificates
```
certbot certonly --manual --work-dir cert --config-dir cert --logs-dir cert --preferred-challenges dns
kubectl create secret generic tls-cert-key --from-file=server.pem=cert/live/service.typeahead.com/fullchain.pem --from-file=privkey.pem=cert/live/service.typeahead.com/privkey.pem
```

# https://medium.com/google-cloud/code-cooking-kubernetes-e715728a578c
# kubectl expose deployment autocomplete --target-port=80 --type=NodePort
# kubectl create -f ingress.yaml
# kubectl get ing --watch
# BACKEND=$(kubectl get ing autocomplete-ingress -o json | jq -j '.metadata.annotations."ingress.kubernetes.io/backends"' | jq -j 'keys[0]')
# gcloud compute backend-services update $BACKEND --enable-cdn
