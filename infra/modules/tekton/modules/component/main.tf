resource "null_resource" "kustomize_tekton" {
  provisioner "local-exec" {
    command = <<EOT
      kustomize build ${path.module} | kubectl apply -f -
    EOT
  }
}
