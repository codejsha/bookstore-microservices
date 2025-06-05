variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "service_accounts" {
  description = "Service account names"
  type = list(string)
}

variable "tekton_address" {
  description = "Tekton address"
  type        = string
}

variable "tekton_fqdn" {
  description = "Tekton FQDN"
  type        = string
}

variable "name_prefix" {
  description = "Resource name prefix"
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
  description = "Harbor address"
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

variable "vault_url" {
  description = "Address"
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
