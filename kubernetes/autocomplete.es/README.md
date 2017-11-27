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
