output "vault_tls_key" {
  value = kubernetes_secret.vault_ha_tls.data["tls.key"]
}

output "vault_tls_crt" {
  value = kubernetes_secret.vault_ha_tls.data["tls.crt"]
}

output "vault_kube_ca_crt" {
  value = kubernetes_secret.vault_ha_tls.data["kubernetes-ca.crt"]
}
