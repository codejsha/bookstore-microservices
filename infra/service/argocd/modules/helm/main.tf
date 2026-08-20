terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "argocd" {
  namespace  = var.namespace
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.4.10"
  values = [
    file("${path.module}/values.yaml"),
    yamlencode({
      configs = {
        secret = {
          argocdServerAdminPassword      = var.admin_password_bcrypt
          argocdServerAdminPasswordMtime = var.admin_password_mtime
        }
        ssh = {
          extraHosts = var.ssh_extra_hosts
        }
      }
    }),
  ]
}
