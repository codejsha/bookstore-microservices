terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_mount" "kv" {
  type        = "kv"
  path        = "kv"
  description = "KV secrets"
  options = { version = "2" }
}
