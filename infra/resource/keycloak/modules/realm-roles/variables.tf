variable "realm_id" {
  description = "Keycloak realm id"
  type        = string
}

variable "roles" {
  description = "Descriptions of the application realm roles, keyed by role name (USER, STAFF, MANAGE, SYSTEM)"
  type        = map(string)

  validation {
    condition     = toset(keys(var.roles)) == toset(["USER", "STAFF", "MANAGE", "SYSTEM"])
    error_message = "roles must describe exactly USER, STAFF, MANAGE and SYSTEM."
  }
}

variable "builtin_default_roles" {
  description = "Keycloak built-in roles kept in the realm default-roles composite alongside USER"
  type        = list(string)
  default     = ["offline_access", "uma_authorization", "account/view-profile", "account/manage-account"]
}
