variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "watch_namespaces" {
  description = "List of namespaces to watch for Kafka custom resources."
  type        = list(string)
}
