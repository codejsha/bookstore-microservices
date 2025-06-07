terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "gitea_ssh_keys" {
  name = "gitea-ssh"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "gitea_ssh_keys" {
  role_name = "gitea-ssh-role"
  bound_service_account_names = ["gitea"]
  bound_service_account_namespaces = [var.namespace]
  token_policies = [vault_policy.gitea_ssh_keys.name]
  token_ttl = "3600" # 1 hour
}
