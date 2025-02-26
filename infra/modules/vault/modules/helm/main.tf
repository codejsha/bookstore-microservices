terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "vault" {
  namespace  = var.namespace
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  version    = "0.29.1"
  values = [
    file("${path.module}/values.yaml")
  ]
}

resource "null_resource" "vault_unseal" {
  depends_on = [helm_release.vault, data.external.vault_keys]
  provisioner "local-exec" {
    command = <<EOT
    kubectl exec -n vault vault-0 -- \
        vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
    EOT
  }
}

resource "null_resource" "vault_raft_join" {
  depends_on = [null_resource.vault_unseal]
  provisioner "local-exec" {
    command = <<EOT
    kubectl exec -n vault -it vault-1 -- vault operator raft join \
        -address=http://vault-1.vault-internal.vault.svc.cluster.local:8200 \
        -leader-ca-cert="${var.vault_kube_ca_crt}" \
        -leader-client-cert="${var.vault_tls_crt}" \
        -leader-client-key="${var.vault_tls_key}" \
        http://vault-0.vault-internal.vault.svc.cluster.local:8200
    kubectl exec -n vault -it vault-2 -- vault operator raft join \
        -address=http://vault-2.vault-internal.vault.svc.cluster.local:8200 \
        -leader-ca-cert="${var.vault_kube_ca_crt}" \
        -leader-client-cert="${var.vault_tls_crt}" \
        -leader-client-key="${var.vault_tls_key}" \
        http://vault-0.vault-internal.vault.svc.cluster.local:8200
    EOT
  }
}

resource "null_resource" "vault_unseal_other_nodes" {
  depends_on = [null_resource.vault_raft_join]
  provisioner "local-exec" {
    command = <<EOT
    kubectl exec -n vault vault-1 -- \
        vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
    kubectl exec -n vault vault-2 -- \
        vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
    EOT
  }
}
