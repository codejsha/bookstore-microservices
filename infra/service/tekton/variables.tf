variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "tekton_address" {
  description = "Tekton address"
  type        = string
}

variable "tekton_service_name" {
  description = "Tekton in-cluster Service name"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address the Chains controller signs against. The controller runs in the cluster, so it must not go out through the edge hostname"
  type        = string
  default     = "http://vault.vault.svc.cluster.local:8200"
}

variable "cosign_key_name" {
  description = "Vault Transit key Chains signs images and provenance with"
  type        = string
  default     = "cosign-key"
}

variable "chains_vault_role" {
  description = "Vault Kubernetes auth role bound to the Chains controller ServiceAccount"
  type        = string
  default     = "tekton-chains-role"
}

variable "builder_id" {
  description = "builder.id stamped into every provenance predicate; the Kyverno provenance policy asserts this exact value"
  type        = string
}

variable "gitea_internal_url" {
  description = "In-cluster Gitea base URL the git resolver calls the SCM API on"
  type        = string
  default     = "http://gitea-http.gitea.svc.cluster.local:3000"
}

variable "gitea_org" {
  description = "Gitea org holding the Tekton catalogs; used as the resolver's default-org so taskRefs only name the repo"
  type        = string
}

variable "resolver_token_secret" {
  description = "Secret holding the Gitea API token the resolver authenticates with"
  type        = string
  default     = "gitea-resolver-token"
}

variable "ca_bundle_configmap" {
  description = "trust-manager bundle ConfigMap distributed into every namespace"
  type        = string
  default     = "internal-ca"
}

variable "ca_bundle_key" {
  description = "Key inside the CA bundle ConfigMap"
  type        = string
  default     = "ca-bundle.crt"
}

variable "prune_ttl_seconds" {
  description = "Delete a finished PipelineRun this long after it completes"
  type        = number
  default     = 604800
}

variable "prune_successful_history_limit" {
  description = "Successful PipelineRuns to retain"
  type        = number
  default     = 20
}

variable "prune_failed_history_limit" {
  description = "Failed PipelineRuns to retain"
  type        = number
  default     = 100
}
