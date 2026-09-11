terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    random = {
      source = "hashicorp/random"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
}

resource "random_password" "valkey" {
  length  = 24
  special = false
}

resource "random_password" "postgresql" {
  length           = 24
  special          = true
  override_special = "!@#$%^&*"
}

resource "random_password" "admin" {
  length  = 32
  special = false
}

resource "vault_policy" "gitea" {
  name   = "gitea"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "gitea" {
  role_name                        = "gitea-role"
  bound_service_account_names      = ["gitea"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.gitea.name]
  token_ttl                        = "3600"
}

resource "vault_kv_secret_v2" "valkey" {
  name  = "gitea/valkey/credentials"
  mount = "kv"
  data_json = jsonencode({
    password = random_password.valkey.result
  })
}

resource "vault_kv_secret_v2" "postgresql" {
  name  = "gitea/postgresql/credentials"
  mount = "kv"
  data_json = jsonencode({
    username = var.postgresql_username
    password = random_password.postgresql.result
  })
}

resource "vault_kv_secret_v2" "admin" {
  name  = "gitea/admin/credentials"
  mount = "kv"
  data_json = jsonencode({
    username = var.admin_username
    password = random_password.admin.result
    email    = var.admin_email
  })
}

resource "kubernetes_manifest" "vaultauth_gitea" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "gitea"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "gitea-role"
        serviceAccount = "gitea"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_gitea_admin" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "gitea-admin"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = "gitea/admin/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "gitea"
      destination = {
        name   = "gitea-admin-credentials"
        create = true
      }
    }
  }
}

resource "tls_private_key" "ssh_host" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "vault_kv_secret_v2" "ssh_host" {
  name  = "gitea/ssh/host"
  mount = "kv"
  data_json = jsonencode({
    private = tls_private_key.ssh_host.private_key_openssh
    public  = trimspace(tls_private_key.ssh_host.public_key_openssh)
  })
}

resource "kubernetes_secret_v1" "ssh_host" {
  metadata {
    name      = var.ssh_host_key_secret_name
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    "gitea.rsa" = tls_private_key.ssh_host.private_key_openssh
  }
}
