terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_mount" "kv_bookstore" {
  type        = "kv"
  path        = "kv-bookstore"
  description = "KV secrets consumed by the bookstore application services"
  options     = { version = "2" }
}

resource "vault_mount" "kv_infra" {
  type        = "kv"
  path        = "kv-infra"
  description = "KV secrets for platform and infrastructure tooling"
  options     = { version = "2" }
}

resource "vault_mount" "transit" {
  type        = "transit"
  path        = "transit"
  description = "Transit engine (cosign image signing)"
}

resource "vault_mount" "database" {
  type        = "database"
  path        = "database"
  description = "Database secrets engine (dynamic creds + static-role rotation)"
}
