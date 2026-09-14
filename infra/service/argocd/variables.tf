variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "argocd_service_name" {
  description = "ArgoCD in-cluster Service name"
  type        = string
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
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

variable "gitea_ssh_fqdn" {
  description = "Gitea SSH service FQDN the known_hosts entries are keyed on"
  type        = string
}

variable "oidc_issuer" {
  description = "OIDC issuer URL of the Keycloak realm that signs Argo CD logins"
  type        = string
}

variable "ca_configmap_name" {
  description = "ConfigMap holding the internal PKI chain that signs the Keycloak edge certificate"
  type        = string
  default     = "vault-pki-ca"
}

variable "ca_configmap_namespace" {
  description = "Namespace of the internal PKI chain ConfigMap"
  type        = string
  default     = "cert-manager"
}

variable "ca_configmap_key" {
  description = "Key of the PEM chain inside the internal PKI ConfigMap"
  type        = string
  default     = "ca-chain.pem"
}
