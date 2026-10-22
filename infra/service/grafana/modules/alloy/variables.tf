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

variable "pod_logs_tail_from_end" {
  description = "Start tailing pod log files from the end when no stored position exists. Keep true for the first rollout of file-based tailing so existing files are not replayed into Loki, then set false so new containers are read from their first line."
  type        = bool
  default     = true
}
