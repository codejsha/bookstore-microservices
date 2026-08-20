include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/vault",
    "../../service/prometheus-crds",
    "../../service/gateway-api",
    "../../service/grafana",
    "../../service/argocd",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "argo-rollouts"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "argo-rollouts"

  gatewayapi_plugin_url = "https://github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/releases/download/v0.16.0/gatewayapi-plugin-linux-amd64"

  grafana_url = "https://grafana.example.com"
}

terraform {
  source = "${get_repo_root()}/infra//service/argo-rollouts"
}
