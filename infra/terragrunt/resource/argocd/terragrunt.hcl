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
    "../../service/argocd",
    "../../service/vault",
    "../../resource/gitea",
    "../../resource/gitea-repo",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "argocd-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace      = "argocd"
  environment    = "dev"
  argocd_address = "argocd.example.com"
  gitea_ssh_fqdn = "gitea-ssh.gitea.svc.cluster.local"
  org_name       = "example-corp"

  app_repos = [
    "catalog-helm", "customer-helm", "identity-helm", "inventory-helm",
    "order-helm", "payment-helm", "delivery-helm", "notification-helm", "support-helm",
    "web-helm", "admin-helm", "admin-web-helm", "settlement-helm",
  ]
}

terraform {
  source = "${get_repo_root()}/infra//resource/argocd"
}
