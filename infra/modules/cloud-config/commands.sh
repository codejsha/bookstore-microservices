######################################################################

### push image
docker login harbor.example.com
docker image push harbor.example.com/bookstore/config-server:0.1.0

### push image with skopeo
skopeo login harbor.example.com
skopeo login --tls-verify=false harbor.example.com
skopeo copy docker-daemon:harbor.example.com/bookstore/config-server:0.1.0 docker://harbor.example.com/bookstore/config-server:0.1.0
skopeo copy --dest-tls-verify=false docker-daemon:harbor.example.com/bookstore/config-server:0.1.0 docker://harbor.example.com/bookstore/config-server:0.1.0

### push helm chart
docker login harbor.example.com
helm package ./helm
helm push config-server-*.tgz oci://harbor.example.com/bookstore-helm-charts
helm push --plain-http config-server-*.tgz oci://harbor.example.com/bookstore-helm-charts

### deploy helm chart
helm install config-server --create-namespace --namespace cloud-config ./config-server-0.1.0.tgz
helm upgrade --install config-server --create-namespace --namespace cloud-config ./config-server-0.1.0.tgz --values ./terraform/modules/helm/values.yaml

helm install config-server --create-namespace --namespace cloud-config config-server --version 0.1.0
helm upgrade --install config-server --create-namespace --namespace cloud-config config-server --version 0.1.0
