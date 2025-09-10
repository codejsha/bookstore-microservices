output "client_uuid" {
  description = "Internal Keycloak UUID of the bookstore-mobile client"
  value       = keycloak_openid_client.mobile.id
}
