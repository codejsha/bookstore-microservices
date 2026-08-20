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
    "../../service/seaweedfs",
    "../../service/prometheus-crds",
    "../../service/gateway-api",
    "../../service/vault-secrets-operator",
  ]
}

locals {
  common          = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role      = "grafana"
  grafana_address = "grafana.example.com"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "grafana"

  grafana_address            = local.grafana_address
  grafana_service_name       = "grafana"
  alloy_grpc_address         = "alloy-grpc.example.com"
  alloy_http_address         = "alloy-http.example.com"
  alloy_gateway_service_name = "alloy-gateway"

  s3_endpoint         = "http://seaweedfs-s3.seaweedfs.svc.cluster.local:8333"
  s3_region           = "us-east-1"
  loki_s3_bucket      = "loki"
  tempo_s3_bucket     = "tempo"
  pyroscope_s3_bucket = "pyroscope"
}

terraform {
  source = "${get_repo_root()}/infra//service/grafana"

  after_hook "reconcile_admin_password" {
    commands = ["apply"]
    execute = [
      "bash", "-c", <<-EOT
      set -euo pipefail
      export VAULT_ADDR='${local.common.vault_url}'
      VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login \
        role=tf-${local.vault_role} \
        jwt="$(kubectl create token tf-${local.vault_role} -n ${local.common.vault_ns} --duration=20m)")
      export VAULT_TOKEN
      user=$(vault kv get -field=admin_user kv/grafana/admin/credentials)
      pw=$(vault kv get -field=admin_password kv/grafana/admin/credentials)
      code=000
      for _ in $(seq 1 30); do
        code=$(curl -sk -o /dev/null -w '%%{http_code}' -u "$user:$pw" "https://${local.grafana_address}/api/org" || true)
        case "$code" in 200|401) break ;; esac
        sleep 5
      done
      if [ "$code" = "401" ]; then
        echo "grafana admin auth drifted from vault — resetting via grafana-cli"
        kubectl exec -n grafana grafana-0 -- grafana-cli admin reset-admin-password "$pw"
      elif [ "$code" != "200" ]; then
        echo "WARN: grafana not reachable (last http $code); skipping admin-password reconcile" >&2
      fi
    EOT
    ]
  }
}
