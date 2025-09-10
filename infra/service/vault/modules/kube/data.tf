data "kubernetes_secret_v1" "vault_token" {
  metadata {
    name      = "vault-token"
    namespace = var.namespace
  }
}
