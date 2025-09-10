resource "null_resource" "kustomize_build" {
  triggers = {
    kustomization = filemd5("${path.module}/kustomization.yaml")
  }
  provisioner "local-exec" {
    command = <<EOT
      kustomize build ${path.module} | kubectl apply -f -
    EOT
  }
}
