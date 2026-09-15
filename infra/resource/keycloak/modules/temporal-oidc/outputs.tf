output "permission_role_ids" {
  description = "IDs of the Temporal client roles, keyed by permission (namespace:role)"
  value       = { for name, role in keycloak_role.permission : name => role.id }
}
