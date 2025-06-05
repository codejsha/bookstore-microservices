data "kubernetes_secret" "vault_tls" {
  metadata {
    name      = "vault-ha-tls"
    namespace = "vault"
  }
}

data "external" "vault_keys" {
  program = ["sh", "${path.module}/cluster-keys.sh"]
}
