terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
  }
}

resource "null_resource" "kustomize_build" {
  triggers = {
    kustomization = filesha256("${path.module}/kustomization.yaml")
  }

  provisioner "local-exec" {
    command = <<EOT
      kustomize build ${path.module} | kubectl apply -f -
    EOT
  }
}
