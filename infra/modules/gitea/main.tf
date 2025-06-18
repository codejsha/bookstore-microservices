terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.37.1"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "= 2.17.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5.0.0"
    }
  }
}

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

resource "kubernetes_namespace" "gitea" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source          = "./modules/helm"
  namespace       = kubernetes_namespace.gitea.metadata[0].name
  admin_email     = var.admin_email
  admin_username  = var.admin_username
  admin_password  = var.admin_password
  valkey_password = var.admin_password
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.gitea.metadata[0].name
  host_address = var.gitea_address
  host_fqdn    = var.gitea_fqdn
  dest_port    = 3000
  name_prefix  = "gitea"
}

module "cert" {
  source           = "./modules/cert"
  namespace        = kubernetes_namespace.gitea.metadata[0].name
  gitea_address    = var.gitea_address
  kube_ca_crt_path = var.kube_ca_crt_path
  providers = {
    vault = vault
  }
}
