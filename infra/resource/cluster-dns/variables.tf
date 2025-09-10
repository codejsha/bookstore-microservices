variable "edge_zone" {
  description = <<-EOT
    DNS zone served by the edge gateway. Every name under it is synthesized to
    the gateway address, so in-cluster clients resolve the same hostnames as
    workstations do instead of each caller needing a Service FQDN.

    A wildcard rather than a hostname list: the zone is ours end to end, it
    already fails to resolve upstream, and an enumerated list silently rots
    every time a route is added. The cost is that a typo resolves instead of
    returning NXDOMAIN: on port 80 the wildcard listener answers 404, while on
    443 every hostname has its own listener and certificate, so an unmatched
    SNI has nothing to terminate the handshake and the connection is reset.
  EOT
  type        = string
  default     = "example.com"
}

variable "gateway_service_name" {
  description = "Istio Gateway API Service the zone resolves to"
  type        = string
  default     = "bookstore-gateway-istio"
}

variable "gateway_namespace" {
  description = "Namespace holding gateway_service_name"
  type        = string
  default     = "istio-system"
}

variable "gateway_address" {
  description = <<-EOT
    Which of the gateway Service's addresses to answer with. "load_balancer"
    hands back the same IP external clients use, so in-cluster and out-of-cluster
    resolution agree; "cluster_ip" keeps the traffic on the cluster network
    instead of hairpinning through the node.
  EOT
  type        = string
  default     = "load_balancer"

  validation {
    condition     = contains(["load_balancer", "cluster_ip"], var.gateway_address)
    error_message = "gateway_address must be load_balancer or cluster_ip."
  }
}

variable "record_ttl" {
  description = "TTL on the synthesized records — short so a gateway address change is picked up without a CoreDNS restart"
  type        = number
  default     = 60
}
