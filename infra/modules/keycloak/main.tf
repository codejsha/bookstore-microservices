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

resource "kubernetes_namespace" "keycloak" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "realm" {
  source    = "./modules/realm"
  namespace = kubernetes_namespace.keycloak.metadata.0.name
  providers = {
    vault = vault
  }
}

module "helm" {
  source         = "./modules/helm"
  namespace      = kubernetes_namespace.keycloak.metadata.0.name
  admin_password = var.admin_username
  admin_username = var.admin_password
  providers = {
    helm = helm
  }
}

module "cert" {
  source           = "./modules/cert"
  namespace        = kubernetes_namespace.keycloak.metadata.0.name
  keycloak_address = var.keycloak_address
  providers = {
    vault = vault
  }
}

module "istio" {
  source    = "./modules/istio"
  namespace = kubernetes_namespace.keycloak.metadata.0.name
}
