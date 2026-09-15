variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "identity_admin_username" {
  description = "Realm-admin service user the identity service authenticates with for Keycloak user administration"
  type        = string
  default     = "identity-admin"
}

variable "identity_admin_email" {
  description = "Email of the identity service realm-admin user"
  type        = string
  default     = "identity-admin@example.com"
}
