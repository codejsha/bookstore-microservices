variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "accounts" {
  description = "Terraform-managed admin console accounts, keyed by a stable account key"
  type = map(object({
    username   = string
    email      = string
    first_name = string
    last_name  = string
    role       = string
  }))
}

variable "role_ids" {
  description = "IDs of the bookstore realm roles, keyed by role name"
  type        = map(string)
}
