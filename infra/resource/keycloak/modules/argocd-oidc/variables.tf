variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "argocd_url" {
  description = "ArgoCD external URL for redirect URIs"
  type        = string
}

variable "namespace" {
  description = "ArgoCD namespace for Kubernetes secret"
  type        = string
}

variable "builtin_default_scopes" {
  description = "Keycloak built-in client scopes kept as defaults on the Argo CD client alongside groups"
  type        = list(string)
  default     = ["profile", "email", "roles", "web-origins", "acr", "basic"]
}

variable "groups_scope_name" {
  description = "Realm client scope that carries platform roles in the groups claim"
  type        = string
}
