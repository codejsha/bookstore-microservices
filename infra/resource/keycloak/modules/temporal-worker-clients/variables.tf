variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "services" {
  description = "Services running Temporal workers; each gets a temporal-worker-<service> client and a kv-bookstore/<service>/temporal secret"
  type        = list(string)
}

variable "permissions" {
  description = "Temporal server permissions (namespace:role) hardcoded into every worker access token"
  type        = list(string)
  default     = ["bookstore:write"]
}

variable "permissions_claim_name" {
  description = "Access token claim carrying the permissions; must match the server's permissionsClaimName"
  type        = string
  default     = "permissions"
}

variable "audience" {
  description = "Audience stamped into worker access tokens; must match the Temporal server's authorization audience"
  type        = string
  default     = "temporal"
}
