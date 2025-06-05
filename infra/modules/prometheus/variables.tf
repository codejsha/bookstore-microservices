variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "prometheus_address" {
  description = "Prometheus address"
  type        = string
}

variable "prometheus_fqdn" {
  description = "Prometheus FQDN"
  type        = string
}

variable "grafana_address" {
  description = "Grafana address"
  type        = string
}

variable "grafana_fqdn" {
  description = "Grafana FQDN"
  type        = string
}
