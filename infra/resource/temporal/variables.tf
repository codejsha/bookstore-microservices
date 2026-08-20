variable "namespace" {
  description = "Kubernetes namespace where Temporal (and the admin-tools image) run."
  type        = string
  default     = "temporal"
}

variable "temporal_address" {
  description = "Temporal frontend gRPC address the admin-tools CLI connects to."
  type        = string
  default     = "temporal-frontend.temporal.svc.cluster.local:7233"
}

variable "temporal_namespaces" {
  description = "Temporal namespaces (NOT Kubernetes namespaces) to register for the bookstore apps."
  type        = list(string)
  default     = ["bookstore"]
}

variable "retention" {
  description = "Retention period for created Temporal namespaces."
  type        = string
  default     = "72h"
}

variable "admintools_image" {
  description = "Temporal admin-tools image providing the `temporal` CLI."
  type        = string
  default     = "temporalio/admin-tools:1.30.3"
}
