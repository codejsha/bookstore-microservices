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

variable "edge_per_ip_ratelimit_enabled" {
  description = "Add per-client-IP token buckets (Envoy local_ratelimit dynamic descriptors keyed on remote_address) next to the host/path-wide edge buckets. Enable only after the gateway access log shows real client addresses; behind a masquerading LB every client shares one bucket"
  type        = bool
  default     = false
}

variable "edge_per_ip_max_tracked_addresses" {
  description = "Upper bound of distinct client addresses each vhost keeps a dynamic bucket for (LRU beyond that)"
  type        = number
  default     = 10000
}

variable "edge_blocked_cidrs" {
  description = "Client CIDRs denied at bookstore-gateway for every host and path (403 before routing). Empty list = no policy. Only meaningful once the gateway sees real client addresses (externalTrafficPolicy Local verified)"
  type        = list(string)
  default     = []
}
