variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "services" {
  description = "List of Python service names"
  type        = list(string)
}

variable "vault_url" {
  description = "Vault URL reachable from outside the cluster"
  type        = string
}

variable "harbor_registry" {
  description = "In-cluster Harbor registry endpoint used as the kaniko push destination"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address used by pipeline steps"
  type        = string
}

variable "nexus_pypi_url" {
  description = "Nexus pypi group simple-index URL the uv builds resolve against (cached pypi.org)"
  type        = string
}
