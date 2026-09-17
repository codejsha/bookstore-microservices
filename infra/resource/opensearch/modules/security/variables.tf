variable "keycloak_issuer_url" {
  description = "Keycloak platform-infra realm issuer URL as published in its OIDC discovery document"
  type        = string
}

variable "oidc_client_id" {
  description = "Keycloak client ID of OpenSearch Dashboards, required as the ID token audience"
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
