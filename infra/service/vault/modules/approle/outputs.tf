output "credentials" {
  description = "AppRole credentials keyed by role name (= each.key in vault_approle_auth_backend_role.service)"
  value = {
    for name, _ in vault_approle_auth_backend_role.service : name => {
      role_id   = data.vault_approle_auth_backend_role_id.service[name].role_id
      secret_id = vault_approle_auth_backend_role_secret_id.service[name].secret_id
    }
  }
  sensitive = true
}

output "role_policies" {
  description = "Map of role name => vault policy name for every approle consumer."
  value       = { for name in keys(local.all_policies) : name => vault_policy.service[name].name }
}
