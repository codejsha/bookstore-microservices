variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kiali_address" {
  description = "Kiali address (Gateway API hostname)"
  type        = string
}

variable "kiali_service_name" {
  description = "Kiali in-cluster Service name"
  type        = string
}

variable "istio_version" {
  description = "Istio chart version for base/istiod/cni/ztunnel (lockstep)"
  type        = string
}

variable "kiali_chart_version" {
  description = "kiali-server helm chart version. Kiali tests each release against the currently-supported Istio releases; 2.26.0 covers Istio 1.29. See https://kiali.io/docs/installation/installation-guide/prerequisites/"
  type        = string
}

variable "edge_hosts" {
  description = "Edge hostnames terminated by per-host HTTPS listeners: hostname => namespaces allowed to bind routes to that listener"
  type        = map(list(string))
  default     = {}
}

variable "grafana_username" {
  description = "Grafana username (kiali data-source)"
  type        = string
  sensitive   = true
}

variable "grafana_password" {
  description = "Grafana password (kiali data-source)"
  type        = string
  sensitive   = true
}

variable "keycloak_realm" {
  description = "Keycloak realm whose token/login-actions paths get dedicated edge rate-limit buckets"
  type        = string
  default     = "bookstore"
}
