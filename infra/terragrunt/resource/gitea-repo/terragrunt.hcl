include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/vault",
    "../../resource/gitea",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  gitea      = read_terragrunt_config(find_in_parent_folders("_envcommon/gitea.hcl")).locals
  vault_role = "gitea-repo-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  org_name       = local.gitea.gitea_org
  gitea_api_url  = "https://git.example.com/api/v1"
  gitea_http_url = "https://git.example.com"

  cloud_config_repo = "cloud-config"
  cloud_config_dir  = "${get_repo_root()}/deploy/kubernetes/cloud-config/repo"

  helm_dir = "${get_repo_root()}/deploy/kubernetes/helm"
  app_repos = [
    "catalog-helm", "customer-helm", "identity-helm", "inventory-helm",
    "order-helm", "payment-helm", "delivery-helm", "notification-helm", "support-helm",
    "web-helm", "admin-helm", "admin-web-helm", "settlement-helm",
  ]

  repo_root = get_repo_root()

  source_repo_suffix = local.gitea.source_repo_suffix
  services = [
    "catalog", "customer", "identity", "inventory",
    "order", "payment", "delivery", "notification", "support",
    "web", "admin", "admin-web", "settlement",
  ]

  tekton_catalog_repo     = "tekton-catalog"
  tekton_catalog_upstream = "https://github.com/tektoncd/catalog"

  github_owner = "codejsha"
  lib_repos    = ["shared-library-go", "shared-library-kotlin", "jooq-codegen-plugin"]

  tekton_custom_catalog_repo = "tektoncd-custom-catalog"
  tekton_custom_catalog_dir  = "${get_repo_root()}/infra/resource/tekton/catalog"
}

terraform {
  source = "${get_repo_root()}/infra//resource/gitea-repo"
}
