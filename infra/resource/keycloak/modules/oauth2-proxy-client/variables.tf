variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "client_id" {
  description = "OIDC client_id for the oauth2-proxy edge/BFF confidential client"
  type        = string
  default     = "bookstore-edge"
}

variable "root_url" {
  description = "Root URL of the browser-facing app fronted by oauth2-proxy"
  type        = string
}

variable "valid_redirect_uris" {
  description = "Allowed OAuth2 redirect URIs — must include the oauth2-proxy callback (https://<app-host>/oauth2/callback)"
  type        = list(string)
}

variable "valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs"
  type        = list(string)
}

variable "web_origins" {
  description = "Allowed CORS origins"
  type        = list(string)
}
