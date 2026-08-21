data "kubernetes_secret" "vault_token" {
  metadata {
    name      = "vault-token"
    namespace = var.namespace
  }
}
