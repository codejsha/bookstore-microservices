include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/gitea",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "gitea-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "gitea"
  gitea_url = "https://git.example.com"
  org_name  = "example-corp"

  admin_username = "gitea_admin"

  dev_repos = [
    "admin-source", "admin-web-source", "catalog-source", "customer-source",
    "identity-source", "inventory-source", "order-source", "payment-source",
    "delivery-source", "notification-source", "support-source", "web-source",
    "settlement-source",

    "shared-library-go",

    "shared-library-kotlin",
    "jooq-codegen-plugin",

    "admin-helm", "admin-web-helm", "catalog-helm", "customer-helm", "identity-helm",
    "inventory-helm", "order-helm", "payment-helm", "delivery-helm",
    "notification-helm", "support-helm", "web-helm", "settlement-helm",
  ]
  devops_repos = []

  dev_usernames    = ["devadmin"]
  devops_usernames = ["devopsadmin"]
}

terraform {
  source = "${get_repo_root()}/infra//resource/gitea"
}
