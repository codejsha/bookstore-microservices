variable "namespace" {
  description = "Bookstore namespace name"
  type        = string
  default     = "bookstore"
}

variable "namespace_labels" {
  description = <<-EOT
    Labels applied to the bookstore namespace. The cluster runs istio AMBIENT
    (ztunnel + istio-cni + waypoint), so the namespace is enrolled via
    istio.io/dataplane-mode=ambient, matching every other platform namespace.
    L7 authz for bookstore is enforced by the waypoint + the AuthorizationPolicies
    in resource/bookstore-mesh.
  EOT
  type        = map(string)
  default = {
    "istio.io/dataplane-mode" = "ambient"
  }
}

variable "registry_host" {
  description = "Container registry host the harbor-pull secret authenticates to."
  type        = string
  default     = "harbor.example.com"
}

variable "harbor_pull_secret_name" {
  description = "Name of the dockerconfigjson image pull secret"
  type        = string
  default     = "harbor-pull"
}

variable "harbor_pull_vault_path" {
  description = "Vault kv (v2) path holding the harbor pull identity {username, password}."
  type        = string
  default     = "harbor/users/harbor-devops/credentials"
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
