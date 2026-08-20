variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "repository_ca_file" {
  description = "Path to the CA cert used to verify the harbor OCI registry's TLS cert"
  type        = string
}

variable "chart_version" {
  description = "config-server chart version, 1.0.0-<8-char git tree hash of helm/> so content changes force an upgrade"
  type        = string
}
