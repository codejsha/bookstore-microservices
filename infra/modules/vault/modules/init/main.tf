resource "null_resource" "vault_unseal" {
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
        http://vault-0.vault-internal.vault.svc.cluster.local:8200
    kubectl exec -n vault -it vault-2 -- vault operator raft join \
        -address=http://vault-2.vault-internal.vault.svc.cluster.local:8200 \
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
