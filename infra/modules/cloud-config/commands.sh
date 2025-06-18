######################################################################

### push helm chart
helm package ./helm
helm push --ca-file ./example-int-ca.crt config-server-*.tgz oci://harbor.example.com/bookstore-helm-charts
helm push --plain-http config-server-*.tgz oci://harbor.example.com/bookstore-helm-charts

### deploy helm chart from archive
helm install config-server --create-namespace --namespace cloud-config ./config-server-0.1.0.tgz --values ./terraform/modules/helm/values.yaml
helm upgrade --install config-server --create-namespace --namespace cloud-config ./config-server-0.1.0.tgz --values ./terraform/modules/helm/values.yaml

### deploy helm chart from repository
helm install config-server --create-namespace --namespace cloud-config bookstore/config-server --version 0.1.0 --values ./terraform/modules/helm/values.yaml
helm upgrade --install config-server --create-namespace --namespace cloud-config bookstore/config-server --version 0.1.0 --values ./terraform/modules/helm/values.yaml
