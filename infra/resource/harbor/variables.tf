variable "harbor_url" {
  description = "Harbor URL"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}


variable "harbor_projects" {
  description = "Harbor project configuration with OIDC group members"
  type = map(object({
    project_name = string
    is_public    = bool
    groups = list(object({
      group_name = string
      role       = string
    }))
  }))
}

variable "robot_name" {
  description = "Harbor robot account used by Tekton to push images"
  type        = string
  default     = "tekton-ci"
}

variable "oidc_issuer" {
  description = "OIDC issuer URL of the Keycloak realm that signs Harbor logins"
  type        = string
}
