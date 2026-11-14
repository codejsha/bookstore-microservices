variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "identity_admin_username" {
  description = "Realm-admin service user the identity service authenticates with for Keycloak user administration"
  type        = string
  default     = "bookstore-identity"
}

variable "identity_admin_email" {
  description = "Email of the identity service realm-admin user"
  type        = string
  default     = "bookstore-identity@example.com"
}

variable "service_account_roles" {
  description = "realm-management client roles granted to the identity client service account"
  type        = list(string)
  default     = ["manage-users", "view-users"]
}
