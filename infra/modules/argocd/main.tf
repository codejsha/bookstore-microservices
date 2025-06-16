terraform {
  required_providers {
    argocd = {
      source  = "argoproj-labs/argocd"
      version = "7.3.1"
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

provider "argocd" {
  server_addr = "${var.argocd_address}:80"
  username    = var.argocd_username
  password    = var.argocd_password
  plain_text  = true
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
  argocd_address = "${var.argocd_address}:80"
  admin_username = var.argocd_username
  admin_password = var.argocd_password
}

module "organization" {
  source     = "./modules/organization"
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  org_name   = var.org_name
  gitea_fqdn = var.gitea_fqdn
  app_repos  = var.app_repos
  providers = {
    argocd = argocd
    vault  = vault
  }
}

module "token" {
  source = "./modules/token"
  providers = {
    argocd = argocd
    vault  = vault
  }
}
