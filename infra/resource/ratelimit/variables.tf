variable "namespace" {
  description = "Namespace the rate-limit service and its EnvoyFilter live in (must be the mesh namespace)"
  type        = string
  default     = "bookstore"
}

variable "image" {
  description = "envoyproxy/ratelimit image"
  type        = string
  default     = "docker.io/envoyproxy/ratelimit:master"
}

variable "waypoint_name" {
  description = "Waypoint gateway whose HTTP chain calls the rate-limit service"
  type        = string
  default     = "bookstore-waypoint"
}

variable "valkey_host" {
  description = "Valkey endpoint backing the rate-limit counters"
  type        = string
  default     = "ratelimit-valkey-headless.bookstore.svc.cluster.local:6379"
}

variable "user_requests_per_second" {
  description = "Per-user (x-user-id) request budget per second across all backend services"
  type        = number
  default     = 30
}

variable "replicas" {
  type    = number
  default = 2
}
