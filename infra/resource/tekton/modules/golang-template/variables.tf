variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_token" {
  description = "ArgoCD token"
  type        = string
  sensitive   = true
}

variable "harbor_address" {
  description = "Harbor address"
  type        = string
}

variable "harbor_username" {
  description = "Harbor username"
  type        = string
  sensitive   = true
}

variable "harbor_token" {
  description = "Harbor token"
  type        = string
  sensitive   = true
}

variable "nexus_docker_registry" {
  description = <<-EOT
    host:port of the Nexus docker connector. Goes into dockerconfig-secret so
    trivy pulls its vulnerability DBs through the ghcr.io cache with the right
    per-host credentials instead of leaking the Nexus login to every registry.
  EOT
  type        = string
}

variable "nexus_username" {
  description = "Nexus username for the docker proxy pull"
  type        = string
}

variable "nexus_password" {
  description = "Nexus password for the docker proxy pull"
  type        = string
  sensitive   = true
}

variable "vault_url" {
  description = "Vault URL reachable from outside the cluster"
  type        = string
}

variable "services" {
  description = "Services built by the golang pipeline. Each needs a Gitea source repo and a <svc>-helm chart repo, both with deploy keys in Vault."
  type        = list(string)
}

variable "lib_repo" {
  description = "Gitea repo holding the shared Go library that services replace into their build"
  type        = string
  default     = "shared-library-go"
}

variable "harbor_registry" {
  description = "In-cluster Harbor registry endpoint used as the kaniko push destination"
  type        = string
}

variable "vault_internal_url" {
  description = "In-cluster Vault address used by pipeline steps"
  type        = string
}

variable "argocd_internal_address" {
  description = "In-cluster ArgoCD API address used by the CI tasks."
  type        = string
  default     = "argocd-server.argocd.svc.cluster.local"
}

variable "nexus_raw_url" {
  description = "Nexus raw repo base (no trailing slash) codegen-verify fetches the mirrored specs from"
  type        = string
}

variable "skip_codegen_verify" {
  description = "\"true\" skips the codegen drift check in CI."
  type        = string
}
