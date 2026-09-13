variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "chart_version" {
  description = "Kyverno Helm chart version"
  type        = string
}

variable "image_pull_secret" {
  description = "dockerconfigjson Secret in this namespace that Kyverno uses to pull signatures and attestations from the private registry. Provisioned by resource/kyverno, not here"
  type        = string
}

variable "ca_bundle_configmap" {
  description = "trust-manager Bundle ConfigMap distributed into every namespace; holds the public roots plus the internal PKI chain"
  type        = string
}

variable "ca_bundle_key" {
  description = "Key inside the CA bundle ConfigMap"
  type        = string
}

variable "vault_url" {
  description = "Vault address for the alert-rule credentials"
  type        = string
}

variable "vault_auth_role" {
  description = "Vault kubernetes-auth role for this unit"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Kubernetes service account JWT used to log in to Vault"
  type        = string
  sensitive   = true
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}
