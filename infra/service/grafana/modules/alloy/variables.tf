variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gateway_target_fqdn" {
  description = "Cluster-local FQDN of the Alloy gateway service (agents forward OTLP here)"
  type        = string
}

variable "loki_push_url" {
  description = "Loki OTLP push URL"
  type        = string
}

variable "tempo_otlp_endpoint" {
  description = "Tempo OTLP gRPC endpoint (host:port)"
  type        = string
}

variable "prometheus_otlp_url" {
  description = "Prometheus OTLP HTTP receiver URL"
  type        = string
}

variable "pyroscope_write_url" {
  description = "Pyroscope push URL (ingester)"
  type        = string
}
