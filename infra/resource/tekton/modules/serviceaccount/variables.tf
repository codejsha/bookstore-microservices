variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "service_account" {
  description = "Service account name"
  type        = string
}

variable "image_pull_secrets" {
  description = "dockerconfigjson Secrets attached to the ServiceAccount. Tekton Chains builds its registry keychain from these, so without them it cannot push signatures or attestations"
  type        = list(string)
  default     = []
}
