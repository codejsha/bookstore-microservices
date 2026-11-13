output "realm_id" {
  description = "Created realm ID"
  value       = keycloak_realm.infra.id
}

output "role_ids" {
  description = "IDs of the platform realm roles, keyed by role name"
  value       = { for name, role in keycloak_role.role : name => role.id }
}

output "groups_scope_name" {
  description = "Client scope that maps platform realm roles into the groups claim"
  value       = keycloak_openid_client_scope.groups.name
}

output "initial_passwords" {
  description = "Initial passwords of the platform accounts that are not stored in Vault, keyed by username"
  value       = { for key, account in var.bootstrap_accounts : account.username => random_password.bootstrap[key].result if !account.store_in_vault }
  sensitive   = true
}
