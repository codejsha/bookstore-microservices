variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "youtrack_address" {
  description = "YouTrack public address"
  type        = string
}

variable "kube_ca_cert" {
  description = "PEM-encoded CA certificate content"
  type        = string
}
