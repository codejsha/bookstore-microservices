data "kubernetes_secret" "vault_tls" {
  metadata {
    name      = "vault-ha-tls"
    namespace = "vault"
  }
}
