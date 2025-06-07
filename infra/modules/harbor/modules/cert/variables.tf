variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "harbor_address" {
  description = "Harbor URL"
  type        = string
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}
