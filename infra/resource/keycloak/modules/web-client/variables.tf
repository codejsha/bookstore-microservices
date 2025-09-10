variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "root_url" {
  description = "Frontend SPA root URL"
  type        = string
}

variable "valid_redirect_uris" {
  description = "Allowed OAuth2 redirect URIs for the SPA"
  type        = list(string)
}

variable "valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs for the SPA"
  type        = list(string)
}

variable "web_origins" {
  description = "Allowed CORS origins for the SPA"
  type        = list(string)
}
