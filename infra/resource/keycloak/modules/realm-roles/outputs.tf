locals {
  roles = {
    USER   = keycloak_role.user
    STAFF  = keycloak_role.staff
    MANAGE = keycloak_role.manage
    SYSTEM = keycloak_role.system
  }
}

output "role_names" {
  description = "Names of the realm roles created for the application"
  value       = [for role in local.roles : role.name]
}

output "role_ids" {
  description = "IDs of the realm roles created for the application, keyed by role name"
  value       = { for name, role in local.roles : name => role.id }
}
