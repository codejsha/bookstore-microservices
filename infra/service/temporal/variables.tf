variable "namespace" {
  description = "Kubernetes namespace for Temporal"
  type        = string
}

variable "temporal_address" {
  description = "Temporal Web UI public address (HTTPRoute hostname)"
  type        = string
}

variable "temporal_authorization_enabled" {
  description = "Enforce JWT authorization (Keycloak JWKS, permissions claim, audience temporal) on the public frontend; enable only after every worker presents tokens"
  type        = bool
  default     = false
}

variable "postgres_cluster_name" {
  description = "CNPG Cluster name (the read-write Service is `<name>-rw` on :5432)"
  type        = string
}

variable "postgres_instances" {
  description = "PostgreSQL instances (1 = single primary, no replica)"
  type        = number
}

variable "postgres_image" {
  description = "CNPG PostgreSQL image (ghcr.io/cloudnative-pg/postgresql:<major>)"
  type        = string
}

variable "postgres_storage_size" {
  description = "Per-instance data volume size"
  type        = string
}

variable "postgres_storage_class_name" {
  description = "StorageClass for the data volume (k3s default is local-path; empty = cluster default)"
  type        = string
}

variable "postgres_resources" {
  description = "Resource requests/limits for the PostgreSQL container"
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
}
