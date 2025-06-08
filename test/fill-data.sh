
curl -X POST --data-binary @is-ip-used.txt http://localhost:8428/api/v1/import/prometheus
curl -X POST --data-binary @metrics.txt http://localhost:8428/api/v1/import/prometheus
curl -X POST --data-binary @registry-optimizer.txt http://localhost:8428/api/v1/import/prometheus
curl -X POST --data-binary @storage-optimizer.txt http://localhost:8428/api/v1/import/prometheus
