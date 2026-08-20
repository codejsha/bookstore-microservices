output "client_uuid" {
  description = "Internal Keycloak UUID of the identity confidential client"
  value       = keycloak_openid_client.identity.id
}
