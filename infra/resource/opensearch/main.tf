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
  }
}

ephemeral "vault_kv_secret_v2" "opensearch" {
  mount = "kv"
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

module "index_template" {
  source = "./modules/index-template"
  providers = {
    opensearch = opensearch
  }
}
