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
