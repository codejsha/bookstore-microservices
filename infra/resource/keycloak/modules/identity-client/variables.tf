variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "realm_admin_username" {
  description = "Realm admin username for identity service"
  type        = string
}

variable "admin_role_id" {
  description = "ID of the ADMIN realm role granted to the realm admin user"
  type        = string
}

variable "realm_admin_first_name" {
  description = "Realm admin first name"
  type        = string
  default     = "DevOps"
}

variable "realm_admin_last_name" {
  description = "Realm admin last name"
  type        = string
  default     = "Admin"
}
