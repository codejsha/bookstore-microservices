variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "client_principals" {
  description = "SPIFFE principals allowed to reach config-server. The config repo is unauthenticated at the app layer (no Spring Security), so ztunnel enforces the allowlist: the backend services that import their config plus the edge gateway serving the oauth2-proxy-gated hostname. Empty disables the policy."
  type        = list(string)
  default     = []
}

variable "repo_root" {
  description = "Absolute path to the monorepo root, used to resolve the helm/ chart tree hash for the chart version"
  type        = string
}
