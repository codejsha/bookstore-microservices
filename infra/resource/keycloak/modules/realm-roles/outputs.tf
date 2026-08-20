output "role_names" {
  description = "Names of the realm roles created for the application"
  value       = [for role in keycloak_role.app : role.name]
}

output "role_ids" {
  description = "IDs of the realm roles created for the application, keyed by role name"
  value       = { for name, role in keycloak_role.app : name => role.id }
}
