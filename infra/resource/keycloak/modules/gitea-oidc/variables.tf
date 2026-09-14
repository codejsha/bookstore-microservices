variable "realm_id" {
  description = "Keycloak realm ID"
  type        = string
}

variable "gitea_url" {
  description = "Gitea external URL for redirect URIs"
  type        = string
}

variable "namespace" {
  description = "Gitea namespace for the OIDC secret and CA bundle"
  type        = string
}

variable "auth_source_name" {
  description = "Name of the Gitea OAuth2 authentication source; forms the callback path"
  type        = string
  default     = "Keycloak"
}

variable "groups_scope_name" {
  description = "Realm client scope that carries platform roles in the groups claim"
  type        = string
}

variable "builtin_default_scopes" {
  description = "Keycloak built-in client scopes kept as defaults on the Gitea client alongside groups"
  type        = list(string)
  default     = ["profile", "email", "roles", "web-origins", "acr", "basic"]
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
