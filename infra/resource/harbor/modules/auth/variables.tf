variable "oidc_endpoint" {
  description = "OIDC issuer URL of the Keycloak realm that signs Harbor logins"
  type        = string
}

variable "oidc_client_id" {
  description = "OIDC client id registered for Harbor"
  type        = string
}

variable "oidc_client_secret" {
  description = "OIDC client secret registered for Harbor"
  type        = string
  sensitive   = true
}

variable "oidc_admin_group" {
  description = "Group claim value whose members become Harbor system administrators"
  type        = string
  default     = "ADMIN"
}
