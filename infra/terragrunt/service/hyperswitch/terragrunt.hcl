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
    "../../service/vault-secrets-operator",
    "../../service/istio",
    "../../service/gateway-api",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "hyperswitch"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "hyperswitch"

  request_timeout = "15s"
}

terraform {
  source = "${get_repo_root()}/infra//service/hyperswitch"

  before_hook "ensure_secrets" {
    commands = ["apply"]
    execute = ["bash", "-c", <<-EOT
      set -euo pipefail
      export VAULT_ADDR='${local.common.vault_url}'
      VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login \
        role=tf-${local.vault_role} \
        jwt="$(kubectl create token tf-${local.vault_role} -n ${local.common.vault_ns} --duration=20m)")
      export VAULT_TOKEN

      if ! vault kv get -field=admin_api_key kv/hyperswitch/app/credentials >/dev/null 2>&1; then
        admin_api_key=$(vault read -field=password sys/policies/password/password-alphanumeric/generate)
        jwt_secret=$(vault read -field=password sys/policies/password/password-alphanumeric/generate)
        recon_admin_api_key=$(vault read -field=password sys/policies/password/password-alphanumeric/generate)
        master_enc_key=$(openssl rand -hex 32)
        user_auth_encryption_key=$(openssl rand -hex 32)
        vault kv put kv/hyperswitch/app/credentials \
          admin_api_key="$admin_api_key" \
          jwt_secret="$jwt_secret" \
          recon_admin_api_key="$recon_admin_api_key" \
          master_enc_key="$master_enc_key" \
          user_auth_encryption_key="$user_auth_encryption_key" \
          >/dev/null
      fi

      if ! vault kv get -field=master_key kv/hyperswitch/card-vault/credentials >/dev/null 2>&1; then
        card_vault_master_key=$(openssl rand -hex 64)
        vault kv put kv/hyperswitch/card-vault/credentials \
          master_key="$card_vault_master_key" \
          >/dev/null
      fi
    EOT
    ]
  }
}
