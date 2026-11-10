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

  harbor_url  = "https://harbor.example.com"
  oidc_issuer = "https://keycloak.example.com/realms/platform-infra"

  harbor_projects = {
    bookstore = {
      project_name = "bookstore"
      is_public    = false
      groups = [
        { group_name = "MANAGER", role = "maintainer" },
        { group_name = "DEVELOPER", role = "developer" },
        { group_name = "PLATFORM_DEVOPS", role = "maintainer" },
        { group_name = "PLATFORM_DEV", role = "developer" },
        { group_name = "BOOKSTORE_ADMIN", role = "projectadmin" },
        { group_name = "BOOKSTORE_DEVOPS", role = "maintainer" },
        { group_name = "BOOKSTORE_DEV", role = "developer" },
      ]
    }
    bookstore-helm-charts = {
      project_name = "bookstore-helm-charts"
      is_public    = false
      groups = [
        { group_name = "MANAGER", role = "maintainer" },
        { group_name = "DEVELOPER", role = "developer" },
        { group_name = "PLATFORM_DEVOPS", role = "maintainer" },
        { group_name = "PLATFORM_DEV", role = "developer" },
        { group_name = "BOOKSTORE_ADMIN", role = "projectadmin" },
        { group_name = "BOOKSTORE_DEVOPS", role = "maintainer" },
        { group_name = "BOOKSTORE_DEV", role = "developer" },
      ]
    }
  }
}

terraform {
  source = "${get_repo_root()}/infra//resource/harbor"
}
