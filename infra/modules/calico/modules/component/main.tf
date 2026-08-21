resource "null_resource" "kustomize_build" {
  provisioner "local-exec" {
    command = <<EOT
      kustomize build ${path.module} | kubectl apply -f -
    EOT
  }
}
