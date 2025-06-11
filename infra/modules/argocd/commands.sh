######################################################################

### initial password
kubectl exec -it -n argocd $(kubectl get pods --selector app.kubernetes.io/name=argocd-server --output jsonpath='{.items[0].metadata.name}') -c server -- argocd admin initial-password

### delete project
argocd proj delete bookstore
