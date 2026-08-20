variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "mysql_services" {
  description = "MySQL services (managed by Oracle MySQL Operator)"
  type        = list(string)
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "mysql_db_config" {
  description = <<-EOT
    Per-service InnoDBCluster metadata. Key = service name (must match an entry in
    mysql_services). Each object:
      cluster_name  - InnoDBCluster .metadata.name (also the Kubernetes Service name prefix)
      database_name - Application schema bootstrapped by the cluster module
      secret_name   - Name of the Kubernetes Opaque secret holding rootUser/rootPassword
      app_user      - Application MySQL user the Vault static role rotates
      instances     - Number of MySQL server instances in the InnoDBCluster
      version       - MySQL server version
      storage_size  - datadirVolumeClaimTemplate request size
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

variable "mysql_client_image" {
  description = "Container image (with a mysql client) used by the DB/app-user bootstrap Job."
  type        = string
  default     = "mysql:8.4"
}

variable "mysqld_exporter_image" {
  description = "mysqld_exporter image for the InnoDBCluster metrics sidecar."
  type        = string
  default     = "prom/mysqld-exporter:v0.14.0"
}

variable "cross_service_readonly" {
  description = <<-EOT
    Cross-service read-only MySQL users. Each entry provisions a SELECT-only user on
    an *existing* service's InnoDBCluster (its target_service — NOT its own cluster) plus a
    Vault database static role that rotates the user's password like the app users.
  EOT
  type = list(object({
    name           = string
    username       = string
    target_service = string
    tables         = list(string)
  }))
  default = []
}
