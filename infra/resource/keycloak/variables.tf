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

variable "kiali_url" {
  description = "Kiali external URL"
  type        = string
}

variable "kiali_namespace" {
  description = "Kiali namespace"
  type        = string
}

variable "grafana_url" {
  description = "Grafana external URL"
  type        = string
}

variable "grafana_namespace" {
  description = "Grafana namespace"
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

variable "bookstore_manager_username" {
  description = "Terraform-managed MANAGE account in the bookstore realm for the admin console"
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

variable "infra_realm_name" {
  description = "Keycloak realm name for platform operators and infrastructure tools"
  type        = string
}



variable "harbor_url" {
  description = "Harbor external URL"
  type        = string
}

variable "infra_bootstrap_accounts" {
  description = "Terraform-managed break-glass accounts in the infra realm, keyed by their Vault path segment"
  type = map(object({
    username   = string
    email      = string
    first_name = string
    last_name  = string
    role       = string
  }))
}

variable "gitea_url" {
  description = "Gitea external URL"
  type        = string
}

variable "gitea_namespace" {
  description = "Gitea namespace"
  type        = string
}

variable "temporal_url" {
  description = "Temporal Web UI external URL"
  type        = string
}

variable "temporal_namespace" {
  description = "Temporal namespace"
  type        = string
}
