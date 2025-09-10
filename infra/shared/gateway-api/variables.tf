variable "hostname" {
  description = "Hostname for the HTTPRoute"
  type        = string
}

variable "service_name" {
  description = "Backend Service short name (HTTPRoute backendRefs.name — not an FQDN)"
  type        = string
}

variable "service_port" {
  description = "Backend service port"
  type        = number
}

variable "route_namespace" {
  description = "Namespace for the HTTPRoute resource"
  type        = string
}

variable "gateway_name" {
  description = "Name of the parent Gateway resource"
  type        = string
  default     = "bookstore-gateway"
}

variable "gateway_namespace" {
  description = "Namespace of the parent Gateway resource"
  type        = string
  default     = "istio-system"
}

variable "name_prefix" {
  description = "Prefix for resource names"
  type        = string
  default     = ""
}

variable "remove_request_headers" {
  description = "Request header names to strip before forwarding upstream."
  type        = list(string)
  default     = []
}

variable "request_timeout" {
  description = "Per-request deadline stamped on the HTTPRoute rule. Defaults to 15s to match the bookstore-mesh waypoint deadline. Pass null to disable"
  type        = string
  default     = "15s"
  nullable    = true
}
