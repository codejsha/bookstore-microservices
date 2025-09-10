variable "namespace" {
  description = "Bookstore application namespace oauth2-proxy is deployed into (shared with the app workloads)"
  type        = string
}

variable "chart_version" {
  description = "oauth2-proxy Helm chart version"
  type        = string
  default     = "10.7.0"
}

variable "credentials_secret_name" {
  description = "Name of the Kubernetes Secret VSO materializes for oauth2-proxy (client-id/client-secret/cookie-secret)"
  type        = string
  default     = "oauth2-proxy-credentials"
}

variable "vault_kv_path" {
  description = "Vault KV v2 path (relative to the kv mount) holding oauth2-proxy client-id/client-secret/cookie-secret"
  type        = string
  default     = "bookstore/oauth2-proxy/config"
}

variable "oidc_issuer_url" {
  description = "OIDC issuer URL for iss validation + discovery — MUST match the Istio RequestAuthentication issuer and the token's iss claim (cluster-internal Keycloak URL)"
  type        = string
}

variable "login_url" {
  description = "Browser-facing Keycloak authorization endpoint (must be reachable from the end-user browser)"
  type        = string
}

variable "redeem_url" {
  description = "Server-side Keycloak token endpoint (cluster-internal; called by oauth2-proxy)"
  type        = string
}

variable "oidc_jwks_url" {
  description = "Server-side Keycloak JWKS endpoint (cluster-internal)"
  type        = string
}

variable "redirect_url" {
  description = "oauth2-proxy callback URL registered on the Keycloak client. Keep it host-relative (/oauth2/callback) so the callback host is derived per request from X-Forwarded-Host."
  type        = string
}

variable "oidc_scope" {
  description = "OIDC scopes requested — keep openid/email/profile/roles so realm_access.roles + claims stay in the token"
  type        = string
  default     = "openid email profile roles"
}

variable "cookie_domains" {
  description = "Cookie domain(s) for the edge session cookie"
  type        = string
}

variable "whitelist_domain" {
  description = "Domains allowed as post-auth redirect targets"
  type        = string
}

variable "session_store" {
  description = "oauth2-proxy session store: 'cookie' (encrypted cookie, default) or 'redis' (Valkey/redis, for refresh-token rotation / large tokens)"
  type        = string
  default     = "cookie"
}

variable "cookie_refresh" {
  description = "How often oauth2-proxy redeems the refresh token against Keycloak, re-checking the session and picking up new realm roles. Must be < access_token_lifespan. Empty disables refresh (no revocation, no role updates)."
  type        = string
  default     = "1m"
}

variable "redis_url" {
  description = "Redis/Valkey connection URL, used only when session_store = redis"
  type        = string
  default     = ""
}

variable "app_host" {
  description = "Browser-facing app hostname fronted by oauth2-proxy (matches the web HTTPRoute host)"
  type        = string
}

variable "service_name" {
  description = "oauth2-proxy in-cluster Service short name (HTTPRoute backendRef)"
  type        = string
  default     = "oauth2-proxy"
}

variable "service_port" {
  description = "oauth2-proxy Service port"
  type        = number
  default     = 4180
}

variable "admin_app_host" {
  description = "Browser-facing admin console hostname fronted by oauth2-proxy (matches the admin-web HTTPRoute host)"
  type        = string
}

variable "admin_web_service_name" {
  description = "Admin console Service name the CUSTOM ext_authz AuthorizationPolicy targets (admin-web Helm release name). The admin BFF gets no waypoint-level CUSTOM policy (it already carries bookstore-introspect, and one CUSTOM per workload); its /api/v1/admin edge traffic is covered by the gateway-attached admin-api-edge-ext-authz instead, which exchanges the session cookie for a bearer before the waypoint's JWT chain runs."
  type        = string
  default     = "admin-web"
}

variable "protected_hosts" {
  description = "Additional edge hostnames placed behind the oauth2-proxy session gate. Each gets the /oauth2 HTTPRoute attached to its listener plus a gateway-attached CUSTOM ext_authz covering every path except /oauth2/*. Use for infrastructure UIs that ship no authentication of their own (e.g. config-server)."
  type        = list(string)
  default     = []
}

variable "extension_provider_name" {
  description = "Name of the envoyExtAuthzHttp extensionProvider registered in the istiod meshConfig"
  type        = string
  default     = "oauth2-proxy"
}

variable "gateway_name" {
  description = "Parent Gateway name for the /oauth2 HTTPRoute"
  type        = string
  default     = "bookstore-gateway"
}

variable "gateway_namespace" {
  description = "Parent Gateway namespace"
  type        = string
  default     = "istio-system"
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}
