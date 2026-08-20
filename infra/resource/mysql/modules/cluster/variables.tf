variable "namespace" {
  description = "Namespace where the InnoDBCluster and bootstrap Job run."
  type        = string
}

variable "mysql_services" {
  description = "List of MySQL service names. Must match the services wired in resource/mysql."
  type        = list(string)
}

variable "mysql_db_config" {
  description = <<-EOT
    Per-service DB metadata map. Key = service name. Each object must supply:
      cluster_name  - InnoDBCluster .metadata.name
      database_name - Application schema bootstrapped by the Job
      secret_name   - Kubernetes secret holding rootUser/rootPassword
      app_user      - Application MySQL user created by the Job
      instances     - Number of InnoDBCluster server instances
      version       - MySQL server version
      storage_size  - datadir PVC request size
  EOT
  type = map(object({
    cluster_name  = string
    database_name = string
    secret_name   = string
    app_user      = string
    instances     = number
    version       = string
    storage_size  = string
  }))
}

variable "mysqld_exporter_image" {
  description = "mysqld_exporter image for the InnoDBCluster metrics sidecar."
  type        = string
  default     = "prom/mysqld-exporter:v0.14.0"
}

variable "mysql_client_image" {
  description = "Container image (with a mysql client) used by the bootstrap Job."
  type        = string
}

variable "cross_service_readonly" {
  description = <<-EOT
    Cross-service read-only MySQL users to provision on an existing (target_service)
    InnoDBCluster. Each entry runs a bootstrap Job (as root) that creates the user and
    grants SELECT only on the named tables of that cluster's application database.
      name           - reader/consumer service
      username       - MySQL user to create
      target_service - key in mysql_db_config of the cluster to read
      tables         - tables to grant SELECT on
  EOT
  type = list(object({
    name           = string
    username       = string
    target_service = string
    tables         = list(string)
  }))
  default = []
}

variable "mysql_resources" {
  description = <<-EOT
    Resource requests/limits for the MySQL server container, applied via the
    InnoDBCluster podSpec. Without this the server pods inherit the namespace
    LimitRange default limit (512Mi), which OOM-kills MySQL 8.4 during init.
  EOT
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
  default = {
    requests = { cpu = "50m", memory = "768Mi" }
    limits   = { cpu = "1", memory = "2Gi" }
  }
}
