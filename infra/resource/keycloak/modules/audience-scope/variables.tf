variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "audience" {
  description = "The aud claim value stamped into tokens (must match bookstore-mesh keycloak_audiences)"
  type        = string
  default     = "bookstore"
}

variable "clients" {
  description = "Map of label => internal Keycloak client UUID to receive the audience scope as a default scope"
  type        = map(string)
}

variable "service_account_clients" {
  description = "Labels from `clients` with service accounts enabled — Keycloak re-attaches the built-in service_account scope to those, so the exhaustive default-scope list has to keep it"
  type        = list(string)
  default     = []
}
