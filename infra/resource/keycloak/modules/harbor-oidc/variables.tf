variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "harbor_url" {
  description = "Harbor external URL for redirect URIs"
  type        = string
}

variable "groups_scope_name" {
  description = "Realm client scope that carries platform roles in the groups claim"
  type        = string
}

variable "builtin_default_scopes" {
  description = "Keycloak built-in client scopes kept as defaults on the Harbor client alongside groups"
  type        = list(string)
  default     = ["profile", "email", "roles", "web-origins", "acr", "basic"]
}
