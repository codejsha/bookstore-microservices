terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "vault_mount" "transit" {
  path        = var.transit_mount_path
  type        = "transit"
  description = "Auto-unseal KMS for the main bookstore Vault (${var.main_vault_namespace})"
}

resource "vault_transit_secret_backend_key" "autounseal" {
  backend          = vault_mount.transit.path
  name             = var.transit_key_name
  deletion_allowed = false
}

resource "vault_policy" "autounseal" {
  name   = "autounseal"
  policy = <<-EOT
    path "${var.transit_mount_path}/encrypt/${var.transit_key_name}" {
      capabilities = ["update"]
    }
    path "${var.transit_mount_path}/decrypt/${var.transit_key_name}" {
      capabilities = ["update"]
    }
  EOT
}

resource "vault_token" "autounseal" {
  policies     = [vault_policy.autounseal.name]
  display_name = "main-vault-autounseal"
  renewable    = true
  period       = var.token_period
  no_parent    = true
  metadata = {
    purpose = "transit-auto-unseal"
    target  = var.main_vault_namespace
  }
}

resource "kubernetes_secret_v1" "autounseal" {
  metadata {
    name      = var.k8s_secret_name
    namespace = var.main_vault_namespace
  }
  data = {
    VAULT_TOKEN = vault_token.autounseal.client_token
  }
  type = "Opaque"
}
