variable "realm_name" {
  description = "Keycloak realm name for platform operators and infrastructure tools"
  type        = string
}

variable "roles" {
  description = "Descriptions of the platform realm roles, keyed by role name"
  type        = map(string)
}

variable "bootstrap_accounts" {
  description = "Terraform-managed break-glass accounts, keyed by their Vault path segment under keycloak/platform/"
  type = map(object({
    username   = string
    email      = string
    first_name = string
    last_name  = string
    role       = string
  }))
}
