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

resource "kubernetes_namespace" "argocd" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.argocd.metadata[0].name
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.argocd.metadata[0].name
  host_address = var.argocd_address
  host_fqdn    = var.argocd_fqdn
  dest_port    = 80
  name_prefix  = "argocd"
}

module "init" {
  source         = "./modules/init"
  namespace      = kubernetes_namespace.argocd.metadata[0].name
  argocd_address = "${var.argocd_address}:80"
  admin_username = var.argocd_username
  admin_password = var.argocd_password
}
