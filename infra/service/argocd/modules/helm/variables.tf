variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "admin_password_bcrypt" {
  description = "bcrypt hash of the argocd admin password (generated once by terragrunt from a Vault password policy). Written to argocd-secret as admin.password."
  type        = string
  sensitive   = true
}

variable "admin_password_mtime" {
  description = "Stable RFC3339 timestamp for argocd-secret admin.passwordMtime. Must be stable across applies (sourced from Vault) to avoid perpetual diffs and spurious password resets."
  type        = string
}

variable "ssh_extra_hosts" {
  description = "known_hosts entries appended to argocd-ssh-known-hosts-cm. Derived from the live Gitea host key so it survives a PVC recreation."
  type        = string
  default     = ""
}

variable "oidc_issuer" {
  description = "OIDC issuer URL of the Keycloak realm that signs Argo CD logins"
  type        = string
}

variable "oidc_root_ca" {
  description = "PEM chain that signs the Keycloak edge certificate, trusted by Argo CD for OIDC discovery and token exchange"
  type        = string
}
