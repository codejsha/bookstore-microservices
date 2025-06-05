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
  token        = module.init.vault_token
  ca_cert_file = var.kube_ca_crt_path
}

resource "kubernetes_namespace" "vault" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "cert" {
  source    = "./modules/cert"
  namespace = kubernetes_namespace.vault.metadata.0.name
  kube_ca_crt = trimspace(file("${path.module}/${var.kube_ca_crt_path}"))
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.vault.metadata.0.name
  providers = {
    helm = helm
  }
}

module "init" {
  source    = "./modules/init"
  namespace = kubernetes_namespace.vault.metadata.0.name
}

module "pki" {
  source    = "./modules/pki"
  vault_url = var.vault_url
  providers = {
    vault = vault
  }
}

module "kv" {
  source = "./modules/kv"
  providers = {
    vault = vault
  }
}

module "kube" {
  source                  = "./modules/kube"
  kube_api_server_address = var.kube_api_server_address
  kube_ca_crt = trimspace(file("${path.module}/${var.kube_ca_crt_path}"))
  providers = {
    vault = vault
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.vault.metadata.0.name
  host_address = var.vault_address
  host_fqdn    = var.vault_fqdn
  name_prefix  = "vault"
}
