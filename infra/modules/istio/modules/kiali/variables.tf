variable "namespace" {
  description = "Namespace name"
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
