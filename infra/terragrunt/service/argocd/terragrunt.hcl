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
    "../../service/gitea",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "argocd"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace           = "argocd"
  argocd_address      = "argocd.example.com"
  argocd_service_name = "argocd-server"

  gitea_ssh_fqdn  = "gitea-ssh.gitea.svc.cluster.local"

  grafana_url = "https://grafana.example.com"
}

terraform {
  source = "${get_repo_root()}/infra//service/argocd"

  before_hook "ensure_admin_credentials" {
    commands = ["plan", "apply"]
    execute = [
      "bash", "-c", <<-EOT
      set -euo pipefail
      export VAULT_ADDR='${local.common.vault_url}'
      VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login \
        role=tf-${local.vault_role} \
        jwt="$(kubectl create token tf-${local.vault_role} -n ${local.common.vault_ns} --duration=20m)")
      export VAULT_TOKEN
      if ! vault kv get -field=password_bcrypt kv/argocd/admin/credentials >/dev/null 2>&1; then
        pw=$(vault read -field=password sys/policies/password/password-special/generate)
        hash=$(printf '%s' "$pw" | htpasswd -niBC 10 "" | cut -d: -f2 | sed 's/^$2y$/$2a$/')
        mtime=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        vault kv put kv/argocd/admin/credentials username=admin password="$pw" password_bcrypt="$hash" password_mtime="$mtime" >/dev/null
      fi
    EOT
    ]
  }
}
