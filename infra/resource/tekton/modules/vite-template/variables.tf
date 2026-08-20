variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL reachable from outside the cluster"
  type        = string
}

variable "services" {
  description = "Apps built by the vite pipeline. Each needs a Gitea source repo and a <svc>-helm chart repo, both with deploy keys in Vault."
  type        = list(string)
}

variable "harbor_registry" {
  description = "In-cluster Harbor registry endpoint used as the kaniko push destination"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address used by pipeline steps"
  type        = string
}

variable "nexus_npm_url" {
  description = "Nexus npm group URL the vite builds resolve against (hosted @bookstore/* + cached registry.npmjs.org)"
  type        = string
}
