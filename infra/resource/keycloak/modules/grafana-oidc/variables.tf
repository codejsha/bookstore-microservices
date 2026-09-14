variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "grafana_url" {
  description = "Grafana external URL for redirect URIs"
  type        = string
}

variable "namespace" {
  description = "Grafana namespace for the OIDC secret and CA bundle"
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
