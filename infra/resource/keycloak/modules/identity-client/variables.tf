variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "manager_username" {
  description = "Realm admin username for identity service"
  type        = string
}

variable "manage_role_id" {
  description = "ID of the MANAGE realm role granted to the realm admin user"
  type        = string
}

variable "manager_first_name" {
  description = "Realm admin first name"
  type        = string
  default     = "DevOps"
}

variable "manager_last_name" {
  description = "Realm admin last name"
  type        = string
  default     = "Admin"
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
