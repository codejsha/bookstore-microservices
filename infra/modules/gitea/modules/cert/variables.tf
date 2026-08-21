variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_address" {
  description = "Gitea address"
  type        = string
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}
