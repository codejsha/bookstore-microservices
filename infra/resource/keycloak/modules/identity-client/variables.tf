variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "service_account_roles" {
  description = "realm-management client roles granted to the identity client service account"
  type        = list(string)
  default     = ["manage-users", "view-users"]
}
