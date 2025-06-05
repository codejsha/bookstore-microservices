# locals {
#   vault_ca_crt  = data.kubernetes_secret.vault_tls.data["ca.crt"]
#   vault_tls_crt = data.kubernetes_secret.vault_tls.data["tls.crt"]
#   vault_tls_key = data.kubernetes_secret.vault_tls.data["tls.key"]
# }
#
# resource "null_resource" "vault_unseal" {
#   provisioner "local-exec" {
#     command = <<EOT
#     kubectl exec -n vault vault-0 -- \
#         vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
#     EOT
#   }
# }
#
# resource "null_resource" "vault_raft_join" {
#   depends_on = [null_resource.vault_unseal]
#   provisioner "local-exec" {
#     command = <<EOT
#     kubectl exec -n vault -it vault-1 -- vault operator raft join \
#         -address=https://vault-1.vault-internal.vault.svc.cluster.local:8200 \
#         -leader-ca-cert="${local.vault_ca_crt}" \
#         -leader-client-cert="${local.vault_tls_crt}" \
#         -leader-client-key="${local.vault_tls_key}" \
#         https://vault-0.vault-internal.vault.svc.cluster.local:8200
#     kubectl exec -n vault -it vault-2 -- vault operator raft join \
#         -address=https://vault-2.vault-internal.vault.svc.cluster.local:8200 \
#         -leader-ca-cert="${local.vault_ca_crt}" \
#         -leader-client-cert="${local.vault_tls_crt}" \
#         -leader-client-key="${local.vault_tls_key}" \
#         https://vault-0.vault-internal.vault.svc.cluster.local:8200
#     EOT
#   }
# }
#
# resource "null_resource" "vault_unseal_other_nodes" {
#   depends_on = [null_resource.vault_raft_join]
#   provisioner "local-exec" {
#     command = <<EOT
#     kubectl exec -n vault vault-1 -- \
#         vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
#     kubectl exec -n vault vault-2 -- \
#         vault operator unseal ${data.external.vault_keys.result.vault_unseal_key}
#     EOT
#   }
# }
