provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "vault" {
  address      = var.vault_url
  token        = var.vault_token
  ca_cert_file = var.kube_ca_crt_path
}

resource "kubernetes_namespace" "minio" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "secret" {
  source         = "./modules/secret"
  namespace      = kubernetes_namespace.minio.metadata.0.name
  admin_username = var.minio_username
  admin_password = var.minio_password
  providers = {
    vault = vault
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.minio.metadata.0.name
  providers = {
    helm = helm
  }
}

module "istio_api" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.minio.metadata.0.name
  host_address = var.minio_api_address
  host_fqdn    = var.minio_fqdn
  dest_port    = 9000
  name_prefix  = "minio-api"
}

module "istio_ui" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.minio.metadata.0.name
  host_address = var.minio_ui_address
  host_fqdn    = var.minio_fqdn
  dest_port    = 9001
  name_prefix  = "minio-ui"
}
