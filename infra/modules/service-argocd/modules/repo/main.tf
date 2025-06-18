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
  git_ssh_url = "git@${var.gitea_fqdn}"
}

resource "argocd_repository" "service_helm_repos" {
  for_each = toset(var.app_repos)
  name            = each.key
  username        = "git"
  repo            = "${local.git_ssh_url}:${var.org_name}/${each.key}.git"
  ssh_private_key = data.vault_generic_secret.repo_ssh_keys[each.key].data["private"]
  insecure        = true
}
