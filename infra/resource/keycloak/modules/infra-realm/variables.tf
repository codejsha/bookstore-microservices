variable "realm_name" {
  description = "Keycloak realm name for platform operators and infrastructure tools"
  type        = string
}

variable "roles" {
  description = "Descriptions of the platform realm roles, keyed by role name"
  type        = map(string)
}

variable "role_composites" {
  description = "Role IDs (such as tool client roles) each platform realm role includes, keyed by role name"
  type        = map(list(string))
  default     = {}
}

variable "bootstrap_accounts" {
  description = "Terraform-managed break-glass accounts, keyed by a stable account key"
  type = map(object({
    username   = string
    email      = string
    first_name = string
    last_name  = string
    role       = string
  }))
}
