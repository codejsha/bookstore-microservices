variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kiali_address" {
  description = "Kiali address"
  type        = string
}

variable "kiali_fqdn" {
  description = "Kiali FQDN"
  type        = string
}

variable "grafana_username" {
  description = "Grafana username"
  type        = string
  sensitive   = true
}

variable "grafana_password" {
  description = "Grafana password"
  type        = string
  sensitive   = true
}
