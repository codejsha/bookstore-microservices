terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.3"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.9"
    }
  }
}

provider "vault" {
  address = var.vault_url
}

data "kubernetes_config_map_v1" "kube_root_ca" {
  metadata {
    name      = "kube-root-ca.crt"
    namespace = "kube-system"
  }
}

locals {
  kube_ca_crt = data.kubernetes_config_map_v1.kube_root_ca.data["ca.crt"]
}

resource "kubernetes_namespace_v1" "vault" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "tls" {
  source      = "./modules/tls"
  namespace   = kubernetes_namespace_v1.vault.metadata[0].name
  kube_ca_crt = local.kube_ca_crt
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace_v1.vault.metadata[0].name
  providers = {
    helm = helm
  }
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
  namespace               = kubernetes_namespace_v1.vault.metadata[0].name
  kube_api_server_address = var.kube_api_server_address
  kube_ca_crt             = local.kube_ca_crt
  tf_auth_roles           = module.approle.role_policies
  depends_on              = [module.helm]
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "approle" {
  source = "./modules/approle"
  providers = {
    vault = vault
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.vault_address
  service_name    = var.vault_service_name
  service_port    = 8200
  route_namespace = kubernetes_namespace_v1.vault.metadata[0].name
  name_prefix     = "vault"
  request_timeout = "10m"
}

module "password_policy" {
  source = "./modules/password-policy"
  providers = {
    vault = vault
  }
  depends_on = [module.helm]
}

resource "vault_pki_secret_backend_role" "wildcard" {
  name               = "wildcard"
  backend            = "pki_int"
  allowed_domains    = ["example.com"]
  allow_subdomains   = true
  allow_bare_domains = true
  allow_glob_domains = true
  max_ttl            = "31536000" # 8760h
  depends_on         = [module.pki]
}

resource "vault_kubernetes_auth_backend_role" "wildcard_issuer" {
  role_name                        = "wildcard-issuer"
  backend                          = "kubernetes"
  bound_service_account_names      = ["wildcard-issuer"]
  bound_service_account_namespaces = ["istio-system"]
  token_policies                   = ["pki_int"]
  token_ttl                        = "3600"
  depends_on                       = [module.kube, module.pki]
}

resource "random_password" "temporal_mysql_root" {
  length  = 32
  special = false
}

resource "random_password" "temporal_mysql_user" {
  length  = 32
  special = false
}

resource "vault_kv_secret_v2" "temporal_mysql" {
  mount = "kv"
  name  = "temporal/mysql"
  data_json = jsonencode({
    root_password = random_password.temporal_mysql_root.result
    username      = "temporal"
    password      = random_password.temporal_mysql_user.result
  })
  depends_on = [module.kv]
}

resource "random_password" "seaweedfs_s3_access_key" {
  length  = 20
  special = false
}

resource "random_password" "seaweedfs_s3_secret_key" {
  length           = 40
  special          = true
  override_special = "!@#$%^&*"
}

resource "vault_kv_secret_v2" "seaweedfs_s3" {
  mount = "kv"
  name  = "seaweedfs/s3/credentials"
  data_json = jsonencode({
    access_key_id     = random_password.seaweedfs_s3_access_key.result
    secret_access_key = random_password.seaweedfs_s3_secret_key.result
  })
  depends_on = [module.kv]
}
