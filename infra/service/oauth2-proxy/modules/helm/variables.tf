variable "namespace" {
  description = "Namespace to deploy oauth2-proxy into"
  type        = string
}

variable "chart_version" {
  description = "oauth2-proxy Helm chart version"
  type        = string
}

variable "existing_secret" {
  description = "Kubernetes Secret (VSO-managed) with client-id/client-secret/cookie-secret"
  type        = string
}

variable "oidc_issuer_url" {
  description = "OIDC issuer URL (cluster-internal; matches token iss + RequestAuthentication)"
  type        = string
}

variable "login_url" {
  description = "Browser-facing Keycloak authorization endpoint"
  type        = string
}

variable "redeem_url" {
  description = "Cluster-internal Keycloak token endpoint"
  type        = string
}

variable "oidc_jwks_url" {
  description = "Cluster-internal Keycloak JWKS endpoint"
  type        = string
}

variable "redirect_url" {
  description = "oauth2-proxy callback URL"
  type        = string
}

variable "oidc_scope" {
  description = "OIDC scopes requested"
  type        = string
}

variable "cookie_domains" {
  description = "Edge session cookie domain"
  type        = string
}

variable "whitelist_domain" {
  description = "Allowed redirect target domains"
  type        = string
}

variable "cookie_refresh" {
  description = "How often oauth2-proxy redeems the refresh token against Keycloak. Bounds how long a revoked user keeps edge access."
  type        = string
}

variable "session_store" {
  description = "'cookie' or 'redis'"
  type        = string
}

variable "redis_url" {
  description = "Redis/Valkey connection URL (session_store = redis only)"
  type        = string
}

variable "service_port" {
  description = "oauth2-proxy Service port"
  type        = number
}
