output "vault_token" {
  value = data.external.vault_keys.result.vault_root_token
}
