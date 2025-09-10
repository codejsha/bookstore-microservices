variable "namespace" {
  description = "Kubernetes namespace where InnoDBCluster instances run (used to read root-credential secrets)"
  type        = string
}

variable "mysql_services" {
  description = "List of MySQL service names. Must match the services wired in resource/mysql (order/payment/delivery/notification)."
  type        = list(string)
}

variable "mysql_db_config" {
  description = <<-EOT
    Per-service DB metadata map. Key = service name. Each object must supply:
      cluster_name  - InnoDBCluster .metadata.name
      database_name - Application schema name
      secret_name   - Kubernetes secret holding rootUser/rootPassword
      app_user      - Application MySQL user the static role rotates
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

variable "cross_service_readonly" {
  description = <<-EOT
    Cross-service read-only MySQL users. For each entry this module creates a Vault static
    role "{name}-{target_service}-mysql-static" on the target_service's existing connection
    (rotating the pre-created read-only user) and extends that connection's allowed_roles.
    The user itself is created by the cluster module's read-only bootstrap Job.
      name           - reader/consumer service
      username       - MySQL user rotated by the static role
      target_service - key in mysql_db_config whose connection hosts the role
      tables         - informational here; the grants are applied by the cluster module
  EOT
  type = list(object({
    name           = string
    username       = string
    target_service = string
    tables         = list(string)
  }))
  default = []
}
