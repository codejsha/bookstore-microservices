variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "kiali_url" {
  description = "Kiali external URL for redirect URIs"
  type        = string
}

variable "namespace" {
  description = "Kiali namespace for the OIDC secret and CA bundle"
  type        = string
}

variable "groups_scope_name" {
  description = "Realm client scope that carries platform roles in the groups claim the Kubernetes API server maps to RBAC groups"
  type        = string
}

variable "builtin_default_scopes" {
  description = "Keycloak built-in client scopes kept as defaults on the Kiali client alongside groups"
  type        = list(string)
  default     = ["profile", "email", "roles", "web-origins", "acr", "basic"]
}

variable "token_lifespan" {
  description = "Access and ID token lifespan in seconds for the Kiali client; Kiali forwards the ID token to the Kubernetes API server for the whole session"
  type        = string
  default     = "1800"
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
