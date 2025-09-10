variable "realm_id" {
  description = "Keycloak realm id"
  type        = string
}

variable "roles" {
  description = "Application realm roles, keyed by role name"
  type        = map(string)
}
