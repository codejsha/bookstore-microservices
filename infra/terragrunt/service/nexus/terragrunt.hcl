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

locals {
  common         = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role     = "nexus"
  admin_username = "admin"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace          = "nexus"
  nexus_address      = "nexus.example.com"
  nexus_service_name = "nexus"

  grafana_url = "https://grafana.example.com"
}

terraform {
  source = "${get_repo_root()}/infra//service/nexus"

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
      if ! vault kv get -field=password kv/nexus/admin/credentials >/dev/null 2>&1; then
        pw=$(vault read -field=password sys/policies/password/password-special/generate)
        vault kv put kv/nexus/admin/credentials \
          username="${local.admin_username}" \
          password="$pw" \
          >/dev/null
      fi
    EOT
    ]
  }

  after_hook "bootstrap_admin_password" {
    commands = ["apply"]
    execute = [
      "bash", "-c", <<-EOT
      set -euo pipefail
      export VAULT_ADDR='${local.common.vault_url}'
      VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login \
        role=tf-${local.vault_role} \
        jwt="$(kubectl create token tf-${local.vault_role} -n ${local.common.vault_ns} --duration=20m)")
      export VAULT_TOKEN

      NS=nexus
      DESIRED_PW=$(vault kv get -field=password kv/nexus/admin/credentials)

      kubectl -n "$NS" rollout status deploy/nexus --timeout=300s || true
      kubectl -n "$NS" port-forward svc/nexus 18081:8081 >/dev/null 2>&1 &
      PF_PID=$!
      trap 'kill "$PF_PID" 2>/dev/null || true' EXIT
      URL="http://127.0.0.1:18081"

      for i in $(seq 1 60); do
        code=$(curl -s -o /dev/null -w '%%{http_code}' "$URL/service/rest/v1/status" || true)
        [ -n "$code" ] && [ "$code" != "000" ] && break
        sleep 5
      done

      auth_ok() { curl -fsS -o /dev/null -u "admin:$1" "$URL/service/rest/v1/security/users?userId=admin"; }

      if auth_ok "$DESIRED_PW"; then
        echo "Nexus admin password already matches Vault; nothing to do."
      else
        POD=$(kubectl get pods -n "$NS" -o name | grep -m1 nexus | cut -d/ -f2)
        INIT_PW=$(kubectl exec -n "$NS" "$POD" -- cat /nexus-data/admin.password 2>/dev/null || true)
        if [ -z "$INIT_PW" ]; then
          echo "ERROR: Vault admin password rejected and no initial /nexus-data/admin.password present; resolve manually." >&2
          exit 1
        fi
        curl -fsS -u "admin:$INIT_PW" \
          -X PUT "$URL/service/rest/v1/security/users/admin/change-password" \
          -H "Content-Type: text/plain" --data-raw "$DESIRED_PW"
        echo "Nexus admin password bootstrapped to the Vault value."
      fi
    EOT
    ]
  }
}
