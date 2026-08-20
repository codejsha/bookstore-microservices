variable "main_vault_namespace" {
  description = "Namespace of the MAIN Vault (where the K8s Secret with the autounseal token is created)"
  type        = string
}

variable "transit_mount_path" {
  description = "Transit engine mount path on the unsealer Vault"
  type        = string
  default     = "transit"
}

variable "transit_key_name" {
  description = "Transit key name used to wrap the main Vault master key. Must match VAULT_TRANSIT_SEAL_KEY_NAME in values.yaml"
  type        = string
  default     = "autounseal"
}

variable "token_period" {
  description = "Periodic TTL of the auto-unseal token (kept alive via seal renewal)"
  type        = string
  default     = "168h"
}

variable "k8s_secret_name" {
  description = "Name of the K8s Secret holding the autounseal token in the main Vault namespace. Must match values.yaml extraSecretEnvironmentVars"
  type        = string
  default     = "vault-transit-unseal"
}
