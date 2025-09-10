variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kiali_chart_version" {
  description = "kiali-server helm chart version. Kiali tests each release against the currently-supported Istio releases; 2.26.0 covers Istio 1.29. See https://kiali.io/docs/installation/installation-guide/prerequisites/"
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
