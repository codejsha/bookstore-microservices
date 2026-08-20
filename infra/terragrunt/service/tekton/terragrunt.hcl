include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

dependencies {
  paths = [
    "../../service/gateway-api",
  ]
}

locals {
  common = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  gitea  = read_terragrunt_config(find_in_parent_folders("_envcommon/gitea.hcl")).locals
}

inputs = {
  namespace           = "tekton-pipelines"
  tekton_address      = "tekton.example.com"
  tekton_service_name = "tekton-dashboard"

  prune_ttl_seconds              = 604800
  prune_successful_history_limit = 20
  prune_failed_history_limit     = 100

  vault_internal_url = "http://vault.vault.svc.cluster.local:8200"
  cosign_key_name   = "cosign-key"
  chains_vault_role = "tekton-chains-role"
  builder_id        = "https://tekton.dev/chains/bookstore-ci"

  gitea_internal_url    = "http://gitea-http.gitea.svc.cluster.local:3000"
  gitea_org             = local.gitea.gitea_org
  resolver_token_secret = "gitea-resolver-token"
}

terraform {
  source = "${get_repo_root()}/infra//service/tekton"
}
