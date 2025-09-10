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
    "../../service/kyverno",
    "../../service/vault",
    "../../resource/harbor",
    "../../resource/tekton",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "kyverno-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace         = "kyverno"
  image_pull_secret = "harbor-pull"

  harbor_address = "harbor.example.com"
  harbor_user    = "harbor-devops"

  cosign_key_name = "cosign-key"

  signed_namespaces = ["bookstore"]
  image_glob        = "harbor.example.com/bookstore/*"
  validation_action = "Deny"
  admission_enabled = true
  provenance_policy_enabled = false
  rekor_url         = "http://rekor.invalid"

  provenance_predicate_type = "https://slsa.dev/provenance/v1"
  builder_id                = "https://tekton.dev/chains/bookstore-ci"
}

terraform {
  source = "${get_repo_root()}/infra//resource/kyverno"
}
