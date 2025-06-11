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

module "component" {
  source = "./modules/component"
}

module "istio" {
  source       = "./modules/istio"
  namespace    = var.namespace
  host_address = var.tekton_address
  host_fqdn    = var.tekton_fqdn
  dest_port    = 9097
  name_prefix  = "tekton-dashboard"

}

module "service_account" {
  for_each = toset(var.service_accounts)
  source          = "./modules/serviceaccount"
  namespace       = var.namespace
  service_account = each.key
}

module "trigger_binding" {
  source    = "./modules/triggerbinding"
  namespace = var.namespace
}

module "golnag_template" {
  source          = "./modules/golang-template"
  namespace       = var.namespace
  argocd_address  = var.argocd_address
  argocd_token    = var.argocd_token
  gitea_address   = var.gitea_address
  harbor_fqdn     = var.harbor_fqdn
  harbor_username = var.harbor_username
  harbor_token    = var.harbor_token
  providers = {
    vault = vault
  }
}
