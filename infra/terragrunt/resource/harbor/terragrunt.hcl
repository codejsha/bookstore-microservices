include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/harbor",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "harbor-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  harbor_url = "https://harbor.example.com"

  harbor_usernames = ["harbor-devops"]

  harbor_projects = {
    bookstore = {
      project_name = "bookstore"
      is_public    = false
      members = [
        { username = "harbor-devops", role = "maintainer" },
      ]
    }
    bookstore-helm-charts = {
      project_name = "bookstore-helm-charts"
      is_public    = false
      members = [
        { username = "harbor-devops", role = "maintainer" },
      ]
    }
  }
}

terraform {
  source = "${get_repo_root()}/infra//resource/harbor"
}
