variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "otel_collector_address" {
  description = "OpenTelemetry Collector gRPC address"
}

variable "otel_collector_fqdn" {
  description = "OpenTelemetry Collector gRPC FQDN"
}
