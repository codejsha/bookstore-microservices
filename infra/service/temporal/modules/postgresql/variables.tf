variable "namespace" {
  description = "Kubernetes namespace for the Temporal PostgreSQL cluster"
  type        = string
}

variable "cluster_name" {
  description = "CNPG Cluster name (the read-write Service is `<name>-rw` on :5432)"
  type        = string
}

variable "instances" {
  description = "PostgreSQL instances (1 = single primary, no replica)"
  type        = number
}

variable "image" {
  description = "CNPG PostgreSQL image (ghcr.io/cloudnative-pg/postgresql:<major>)"
  type        = string
}

variable "storage_size" {
  description = "Per-instance data volume size"
  type        = string
}

variable "storage_class_name" {
  description = "StorageClass for the data volume (empty = cluster default)"
  type        = string
}

variable "resources" {
  description = "Resource requests/limits for the PostgreSQL container"
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
}
