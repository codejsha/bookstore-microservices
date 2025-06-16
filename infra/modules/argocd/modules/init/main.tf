resource "null_resource" "update_admin_password" {
  provisioner "local-exec" {
    command = <<EOT
      password=$(kubectl exec -n argocd $(kubectl get pods --selector app.kubernetes.io/name=argocd-server --output jsonpath='{.items[0].metadata.name}') -c server \
        -- argocd admin initial-password | grep -m1 -v '^$')
      argocd login --username ${var.admin_username} --password $password --grpc-web --insecure --plaintext ${var.argocd_address}
      argocd account update-password --account ${var.admin_username} --current-password $password --new-password ${var.admin_password} \
        --grpc-web --insecure --plaintext --server ${var.argocd_address}
    EOT
  }
}
