variable "namespace" {
  description = "Kubernetes namespace for the Temporal MySQL cluster"
  type        = string
}

variable "cluster_name" {
  description = "InnoDBCluster name (also the RW Service name on :3306)"
  type        = string
}

variable "root_password" {
  description = "MySQL root password (generated in the parent, stored in Vault)"
  type        = string
  sensitive   = true
}

variable "instances" {
  description = "MySQL server instances (3 for Group Replication quorum; 1 = no HA, cannot self-heal after a restart)"
  type        = number
}

variable "router_instances" {
  description = "MySQL Router instances"
  type        = number
}

variable "mysql_version" {
  description = "MySQL server version"
  type        = string
}

variable "storage_size" {
  description = "Per-instance data volume size"
  type        = string
}

variable "storage_class_name" {
  description = "StorageClass for the data volume (k3s default is local-path)"
  type        = string
}

variable "db_secret_name" {
  description = "Secret holding the password under key `password` for the Temporal chart's existingSecret"
  type        = string
}

variable "mysql_resources" {
  description = <<-EOT
    Resource requests/limits for the MySQL server container, applied via the
    InnoDBCluster podSpec. Without this the server pods fall back to the
    namespace LimitRange default request (10m/32Mi) with no limit, so the
    scheduler packs them as if they were idle.
  EOT
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
  default = {
    requests = { cpu = "100m", memory = "512Mi" }
    limits   = { cpu = "1", memory = "2Gi" }
  }
}

variable "router_resources" {
  description = "Resource requests/limits for the MySQL Router container"
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
  default = {
    requests = { cpu = "50m", memory = "128Mi" }
    limits   = { cpu = "500m", memory = "256Mi" }
  }
}
