variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "argocd_token" {
  description = "ArgoCD token"
  type        = string
  sensitive   = true
}

variable "harbor_fqdn" {
  description = "Harbor FQDN"
  type        = string
}

variable "harbor_username" {
  description = "Harbor username"
  type        = string
}

variable "harbor_token" {
  description = "Harbor token"
  type        = string
  sensitive   = true
}

variable "gitea_address" {
  description = "Host FQDN"
  type        = string
}
