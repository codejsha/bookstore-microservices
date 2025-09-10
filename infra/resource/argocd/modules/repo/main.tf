terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

locals {
  git_ssh_url = "ssh://git@${var.gitea_ssh_fqdn}:${var.gitea_ssh_port}/${var.org_name}"
}

data "vault_kv_secret_v2" "repo_ssh_keys" {
  for_each = toset(var.app_repos)
  mount    = "kv"
  name     = "gitea/ssh/${each.key}"
}

resource "argocd_repository" "service_helm_repos" {
  for_each        = toset(var.app_repos)
  name            = each.key
  repo            = "${local.git_ssh_url}/${each.key}.git"
  ssh_private_key = data.vault_kv_secret_v2.repo_ssh_keys[each.key].data["private"]
}
