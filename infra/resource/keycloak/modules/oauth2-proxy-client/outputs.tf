output "client_id" {
  description = "The oauth2-proxy confidential client_id"
  value       = keycloak_openid_client.oauth2_proxy.client_id
}

output "client_uuid" {
  description = "Internal Keycloak UUID of the oauth2-proxy client"
  value       = keycloak_openid_client.oauth2_proxy.id
}

output "client_secret" {
  description = "The oauth2-proxy confidential client secret (also written to Vault kv/bookstore/oauth2-proxy/config)"
  value       = keycloak_openid_client.oauth2_proxy.client_secret
  sensitive   = true
}

output "vault_kv_path" {
  description = "Vault KV v2 path (relative to the kv mount) holding client-id/client-secret/cookie-secret"
  value       = vault_kv_secret_v2.oauth2_proxy.name
}
