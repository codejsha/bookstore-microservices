variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_url" {
  description = "ArgoCD URL"
  type        = string
}

variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "argocd_fqdn" {
  description = "ArgoCD FQDN"
  type        = string
}

variable "argocd_username" {
  description = "ArgoCD username"
  type        = string
  sensitive   = true
}

variable "argocd_password" {
  description = "ArgoCD password"
  type        = string
  sensitive   = true
}

variable "gitea_fqdn" {
  description = "Gitea FQDN"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}
