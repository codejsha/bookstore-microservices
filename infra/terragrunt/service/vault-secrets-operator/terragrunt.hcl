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
    "../../service/vault",
  ]
}

inputs = {
  namespace = "vault-secrets-operator"

  vault_address = "http://vault.vault.svc.cluster.local:8200"
}

terraform {
  source = "${get_repo_root()}/infra//service/vault-secrets-operator"
}
