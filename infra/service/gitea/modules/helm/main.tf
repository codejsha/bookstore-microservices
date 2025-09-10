terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_service_account_v1" "gitea" {
  metadata {
    name      = "gitea"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "gitea" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "gitea-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.gitea.metadata[0].name
    }
  }
}

resource "helm_release" "gitea" {
  namespace  = var.namespace
  name       = "gitea"
  repository = "https://dl.gitea.com/charts"
  chart      = "gitea"
  version    = "12.7.0"
  values = [
    file("${path.module}/values.yaml"),
    yamlencode({
      extraVolumes = [
        {
          name = "ssh-host-key"
          secret = {
            secretName  = var.ssh_host_key_secret_name
            defaultMode = 416
          }
        },
      ]
      extraContainerVolumeMounts = [
        {
          name      = "ssh-host-key"
          mountPath = "/etc/gitea/ssh"
          readOnly  = true
        },
      ]
    }),
  ]
  timeout = 180
  set_sensitive = [
    { name = "gitea.admin.email", value = var.admin_email },
    { name = "valkey.global.valkey.password", value = var.valkey_password },
    { name = "postgresql.global.postgresql.auth.username", value = var.postgresql_username },
    { name = "postgresql.global.postgresql.auth.password", value = var.postgresql_password }
  ]

  set = [
    { name = "gitea.admin.existingSecret", value = "gitea-admin-credentials" }
  ]
}
