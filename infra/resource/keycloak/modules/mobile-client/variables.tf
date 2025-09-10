variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "valid_redirect_uris" {
  description = "Allowed redirect URIs for the native app (custom scheme matches app.json)"
  type        = list(string)
}

variable "valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs for the native app"
  type        = list(string)
}
