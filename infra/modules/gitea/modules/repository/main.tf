terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "gitea_repository" "git_repo" {
  name        = var.repo_name
  description = var.repo_description
  username    = var.gitea_username
  private     = true
  auto_init   = false
  license     = "Apache-2.0"
}

resource "tls_private_key" "tls_keys" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "gitea_repository_key" "repo_key" {
  title      = "${var.repo_name} repository ssh key"
  repository = gitea_repository.git_repo.id
  key        = tls_private_key.tls_keys.public_key_openssh
}

resource "vault_kv_secret_v2" "tls_keys" {
  name  = "gitea/ssh/${var.repo_name}"
  mount = "kv"
  data_json = jsonencode(
    {
      private = tls_private_key.tls_keys.private_key_pem,
      public  = tls_private_key.tls_keys.public_key_openssh
    }
  )
}
