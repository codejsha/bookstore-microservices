variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kube_ca_crt" {
  description = "Kubernetes CA certificate"
  type        = string
}

variable "local_output_dir" {
  description = "Absolute directory where the Vault TLS key and certificate are written"
  type        = string
}
