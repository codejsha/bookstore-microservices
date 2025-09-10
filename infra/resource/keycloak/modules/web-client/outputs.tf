output "client_uuid" {
  description = "Internal Keycloak UUID of the bookstore-web client"
  value       = keycloak_openid_client.web.id
}
