variable "namespace" {
  description = "Kubernetes namespace where the CNPG clusters live."
  type        = string
}

variable "postgres_services" {
  description = "List of service names whose CNPG clusters should be registered with the Vault database secrets engine. Each entry must match the service name used as the Helm release / clusterName prefix."
  type        = list(string)
}

variable "postgres_db_map" {
  description = <<-EOT
    Map of service name → database name as declared in the CNPG Cluster spec
    (Cluster.spec.bootstrap.initdb.database). Must match values.yaml postgresql.database.
    Example: { catalog = "catalog_db", customer = "customer_db" }
  EOT
  type        = map(string)
}

variable "postgres_app_user_map" {
  description = <<-EOT
    Map of service name → application username as declared in the CNPG Cluster spec
    (Cluster.spec.bootstrap.initdb.owner). Used as the static-role username.
    Example: { catalog = "catalog", customer = "customer" }
  EOT
  type        = map(string)
}

variable "cross_service_readonly" {
  description = <<-EOT
    Cross-service read-only Postgres roles. For each entry this module creates a Vault
    static role "{name}-{target_service}-postgres-static" on the target_service's existing
    connection (rotating the pre-created read-only role) and extends that connection's
    allowed_roles. The role itself is created by the cluster module's read-only bootstrap Job.
      name           - reader/consumer service
      username       - Postgres role rotated by the static role
      target_service - key in postgres_services whose connection hosts the role
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

variable "rotation_schedule" {
  description = "Cron schedule (Vault server clock, UTC) on which the database static roles rotate their passwords. Mutually exclusive with rotation_period."
  type        = string
  default     = "10 19 * * *"
}

variable "rotation_window" {
  description = "Seconds after each rotation_schedule tick during which Vault may still perform the rotation before skipping it."
  type        = number
  default     = 3600
}
