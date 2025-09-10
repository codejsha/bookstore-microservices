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
    "../../service/gateway-api",
    "../../service/cert-manager",
    "../../service/prometheus-crds",
    "../../service/vault-secrets-operator",
    "../../service/grafana",
  ]
}

dependency "vault" {
  config_path = "../../service/vault"

  mock_outputs = {
    seaweedfs_s3_credentials = {
      access_key_id     = "mock-access-key-id"
      secret_access_key = "mock-secret-access-key"
    }
  }
  mock_outputs_allowed_terraform_commands = ["validate", "plan"]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "harbor"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace           = "harbor"
  harbor_address      = "harbor.example.com"
  harbor_service_name = "harbor"
  bucket_names        = ["harbor-storage"]

  grafana_url = "https://grafana.example.com"

  aws_s3_api_url = "https://seaweedfs-api.example.com"
  aws_access_key = dependency.vault.outputs.seaweedfs_s3_credentials.access_key_id
  aws_secret_key = dependency.vault.outputs.seaweedfs_s3_credentials.secret_access_key
}

terraform {
  source = "${get_repo_root()}/infra//service/harbor"

  before_hook "ensure_admin_password" {
    commands = ["apply"]
    execute = [
      "bash", "-c", <<-EOT
      set -euo pipefail
      export VAULT_ADDR='${local.common.vault_url}'
      VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login \
        role=tf-${local.vault_role} \
        jwt="$(kubectl create token tf-${local.vault_role} -n ${local.common.vault_ns} --duration=20m)")
      export VAULT_TOKEN

      if ! vault kv get -field=password kv/harbor/admin/credentials >/dev/null 2>&1; then
        pw=$(vault read -field=password sys/policies/password/password-special/generate)
        vault kv put kv/harbor/admin/credentials username=admin password="$pw" >/dev/null
      fi

      if ! vault kv get -field=REGISTRY_PASSWD kv/harbor/registry/credentials >/dev/null 2>&1; then
        pw=$(vault read -field=password sys/policies/password/password-alphanumeric/generate)
        htp=$(htpasswd -nbB harbor_registry_user "$pw")
        vault kv put kv/harbor/registry/credentials \
          REGISTRY_PASSWD="$pw" \
          REGISTRY_HTPASSWD="$htp" \
          username="harbor_registry_user" >/dev/null
      fi
    EOT
    ]
  }
}
