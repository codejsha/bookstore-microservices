variable "keycloak_url" {
  description = "Keycloak URL"
  type        = string
}

variable "realm_name" {
  description = "Keycloak realm name"
  type        = string
}

variable "keycloak_namespace" {
  description = "Namespace running the Keycloak StatefulSet"
  type        = string
}

variable "master_admin_username" {
  description = "Permanent master-realm admin username that replaces the operator's temporary admin"
  type        = string
}

variable "master_admin_email" {
  description = "Permanent master-realm admin email"
  type        = string
}

variable "argocd_url" {
  description = "ArgoCD external URL"
  type        = string
}

variable "argocd_namespace" {
  description = "ArgoCD namespace"
  type        = string
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

variable "web_root_url" {
  description = "Web SPA root URL"
  type        = string
}

variable "web_valid_redirect_uris" {
  description = "Allowed OAuth2 redirect URIs for the SPA"
  type        = list(string)
}

variable "web_valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs for the SPA"
  type        = list(string)
}

variable "web_web_origins" {
  description = "Allowed CORS origins for the SPA"
  type        = list(string)
}

variable "identity_realm_admin_username" {
  description = "Realm admin username for identity service"
  type        = string
}

variable "mobile_valid_redirect_uris" {
  description = "Allowed redirect URIs for the native app (custom scheme matches app.json)"
  type        = list(string)
}

variable "mobile_valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs for the native app"
  type        = list(string)
}

variable "oauth2_proxy_client_id" {
  description = "OIDC client_id for the oauth2-proxy edge/BFF confidential client"
  type        = string
  default     = "bookstore-edge"
}

variable "oauth2_proxy_root_url" {
  description = "Root URL of the browser-facing app fronted by oauth2-proxy"
  type        = string
}

variable "oauth2_proxy_valid_redirect_uris" {
  description = "Allowed redirect URIs for oauth2-proxy — must include https://<app-host>/oauth2/callback"
  type        = list(string)
}

variable "oauth2_proxy_valid_post_logout_redirect_uris" {
  description = "Allowed post-logout redirect URIs for oauth2-proxy"
  type        = list(string)
}

variable "oauth2_proxy_web_origins" {
  description = "Allowed CORS origins for oauth2-proxy"
  type        = list(string)
}
