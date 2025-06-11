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

resource "kubernetes_namespace" "opensearch" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "secret" {
  source                 = "./modules/secret"
  namespace              = var.namespace
  initial_admin_password = var.initial_admin_password
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.opensearch.metadata.0.name
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.opensearch.metadata.0.name
  host_address = var.opensearch_address
  host_fqdn    = var.opensearch_fqdn
  dest_port    = 5601
  name_prefix  = "opensearch"
}
