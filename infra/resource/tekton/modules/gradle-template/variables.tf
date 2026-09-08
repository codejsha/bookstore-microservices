variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL reachable from outside the cluster"
  type        = string
}

variable "services" {
  description = "Services built by the gradle pipeline. Each needs a Gitea source repo and a <svc>-helm chart repo, both with deploy keys in Vault."
  type        = list(string)
}

variable "harbor_registry" {
  description = "In-cluster Harbor registry endpoint used as the kaniko push destination"
  type        = string
}

variable "harbor_address" {
  description = "External Harbor hostname used for step images the node pulls (with the harbor-pull secret)"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address used by pipeline steps"
  type        = string
}

variable "nexus_maven_url" {
  description = "Nexus maven group URL the Gradle builds resolve against (hosted + cached Central + Gradle Plugin Portal)"
  type        = string
}

variable "lib_repos" {
  description = "Gitea repos holding the sibling Kotlin sources every gradle service includes as a composite build (plugin + shared library)"
  type        = set(string)
  default     = ["jooq-codegen-plugin", "shared-library-kotlin"]
}

variable "nexus_raw_url" {
  description = "Nexus raw-hosted repo base (no trailing slash) codegen-verify-external / codegen-verify-grpc fetch the mirrored specs from"
  type        = string
}

variable "codegen_verify_targets" {
  description = "Codegen drift checks the CI pipeline runs: any of openapi, external, grpc, db. Empty runs none."
  type        = list(string)
}
