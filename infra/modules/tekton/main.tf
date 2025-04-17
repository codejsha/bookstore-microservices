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

resource "kubernetes_namespace" "tekton" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "component" {
  source = "./modules/component"
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.tekton.metadata.0.name
  host_address = var.tekton_address
  host_fqdn    = var.tekton_fqdn
  name_prefix  = var.name_prefix
}

module "service_account" {
  for_each = toset(var.service_accounts)
  source          = "./modules/serviceaccount"
  namespace       = kubernetes_namespace.tekton.metadata.0.name
  service_account = each.key
}

module "trigger_binding" {
  source    = "./modules/triggerbinding"
  namespace = kubernetes_namespace.tekton.metadata.0.name
}

module "golnag_template" {
  source          = "./modules/golang-template"
  namespace       = kubernetes_namespace.tekton.metadata.0.name
  argocd_address  = var.argocd_address
  argocd_token    = var.argocd_token
  gitea_fqdn      = var.gitea_fqdn
  harbor_fqdn     = var.harbor_fqdn
  harbor_username = var.harbor_username
  harbor_token    = var.harbor_token
  providers = {
    vault = vault
  }
}
