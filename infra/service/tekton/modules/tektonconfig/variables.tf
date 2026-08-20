variable "namespace" {
  description = "Target namespace the operator installs Tekton components into"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address the Chains controller signs against"
  type        = string
}

variable "cosign_key_name" {
  description = "Vault Transit key Chains signs images and provenance with"
  type        = string
}

variable "chains_vault_role" {
  description = "Vault Kubernetes auth role bound to the Chains controller ServiceAccount. Created by resource/tekton, so Chains cannot sign until that unit has run"
  type        = string
}

variable "builder_id" {
  description = "builder.id stamped into every provenance predicate; the Kyverno provenance policy asserts this exact value"
  type        = string
}

variable "gitea_internal_url" {
  description = "In-cluster Gitea base URL the git resolver calls the SCM API on"
  type        = string
}

variable "gitea_org" {
  description = "Gitea org holding the Tekton catalogs; used as the resolver's default-org so taskRefs only name the repo"
  type        = string
}

variable "resolver_token_secret" {
  description = "Secret holding the Gitea API token the resolver authenticates with. Created by resource/tekton"
  type        = string
}

variable "ca_bundle_configmap" {
  description = "trust-manager bundle ConfigMap mounted into the Chains controller so it can reach Harbor over the internal PKI"
  type        = string
}

variable "ca_bundle_key" {
  description = "Key inside the CA bundle ConfigMap"
  type        = string
}

variable "prune_ttl_seconds" {
  description = "Delete a finished PipelineRun this long after it completes. Outer bound — the history limits usually bite first"
  type        = number
}

variable "prune_successful_history_limit" {
  description = "Successful PipelineRuns to retain"
  type        = number
}

variable "prune_failed_history_limit" {
  description = "Failed PipelineRuns to retain. Kept well above the successful limit so a broken build is still there to read"
  type        = number
}
