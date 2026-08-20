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
