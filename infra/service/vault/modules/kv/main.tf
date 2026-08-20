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
