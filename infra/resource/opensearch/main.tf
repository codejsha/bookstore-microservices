terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    opensearch = {
      source  = "opensearch-project/opensearch"
      version = "~> 2.4"
    }
    restapi = {
      source  = "Mastercard/restapi"
      version = "~> 3.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

ephemeral "vault_kv_secret_v2" "opensearch" {
  mount = "kv-infra"
  name  = "opensearch/admin/credentials"
}

provider "opensearch" {
  url               = var.opensearch_api_url
  username          = ephemeral.vault_kv_secret_v2.opensearch.data["username"]
  password          = ephemeral.vault_kv_secret_v2.opensearch.data["initial_password"]
  insecure          = true
  healthcheck       = false
  sign_aws_requests = false
}

provider "restapi" {
  uri      = var.opensearch_api_url
  username = ephemeral.vault_kv_secret_v2.opensearch.data["username"]
  password = ephemeral.vault_kv_secret_v2.opensearch.data["initial_password"]
  insecure = true
  headers = {
    "Content-Type" = "application/json"
  }
}

module "index_template" {
  source = "./modules/index-template"
  providers = {
    opensearch = opensearch
  }
}

module "security" {
  source              = "./modules/security"
  keycloak_issuer_url = var.keycloak_issuer_url
  oidc_client_id      = var.oidc_client_id
  providers = {
    opensearch = opensearch
    restapi    = restapi
    kubernetes = kubernetes
  }
}
