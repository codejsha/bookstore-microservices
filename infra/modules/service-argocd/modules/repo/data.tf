data "vault_generic_secret" "repo_ssh_keys" {
  for_each = toset(var.app_repos)
  path = "kv/gitea/ssh/${each.key}"
}
