locals {
  manifest = templatefile("${path.module}/tektonconfig.yaml.tftpl", {
    namespace                      = var.namespace
    prune_ttl_seconds              = var.prune_ttl_seconds
    prune_successful_history_limit = var.prune_successful_history_limit
    prune_failed_history_limit     = var.prune_failed_history_limit
    vault_url                      = var.vault_internal_url
    cosign_key_name                = var.cosign_key_name
    chains_vault_role              = var.chains_vault_role
    builder_id                     = var.builder_id
    gitea_internal_url             = var.gitea_internal_url
    gitea_org                      = var.gitea_org
    resolver_token_secret          = var.resolver_token_secret
    ca_bundle_configmap            = var.ca_bundle_configmap
    ca_bundle_key                  = var.ca_bundle_key
  })
}

resource "local_file" "tekton_config" {
  filename        = "${path.module}/tektonconfig.rendered.yaml"
  content         = local.manifest
  file_permission = "0644"
}

resource "null_resource" "tekton_config" {
  triggers = {
    manifest = sha1(local.manifest)
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      set -euo pipefail
      for attempt in $(seq 1 12); do
        if kubectl apply --server-side --force-conflicts -f ${local_file.tekton_config.filename}; then
          break
        fi
        if [ "$attempt" -eq 12 ]; then
          exit 1
        fi
        sleep 10
      done
      kubectl wait --for=condition=Ready tektonconfig/config --timeout=900s
    EOT
  }
}
