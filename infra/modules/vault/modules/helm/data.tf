data "external" "vault_keys" {
  depends_on = [helm_release.vault]
  program = ["sh", "${path.module}/cluster-keys.sh"]
}
