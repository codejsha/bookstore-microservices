output "client_scope_name" {
  description = "Name of the created audience client scope"
  value       = keycloak_openid_client_scope.bookstore_audience.name
}
