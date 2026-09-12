variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "realm_admin_username" {
  description = "Realm admin username for identity service"
  type        = string
}

variable "manage_role_id" {
  description = "ID of the MANAGE realm role granted to the realm admin user"
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
