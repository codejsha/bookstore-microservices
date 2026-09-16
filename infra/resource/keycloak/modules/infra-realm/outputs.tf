output "realm_id" {
  description = "Created realm ID"
  value       = keycloak_realm.infra.id
}

output "role_ids" {
  description = "IDs of the platform realm roles, keyed by role name"
  value = {
    DEVELOPER = keycloak_role.developer.id
    MANAGER   = keycloak_role.manager.id
    ADMIN     = keycloak_role.admin.id
  }
}

output "groups_scope_name" {
  description = "Client scope that maps platform realm roles into the groups claim"
  value       = keycloak_openid_client_scope.groups.name
}
