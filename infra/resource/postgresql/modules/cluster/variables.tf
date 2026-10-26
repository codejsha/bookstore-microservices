variable "namespace" {
  description = "Kubernetes namespace where the CloudNativePG Cluster CRs are created."
  type        = string
}

variable "clusters" {
  description = <<-EOT
    Per-service CloudNativePG Cluster spec, keyed by service name.
    Each object mirrors the values previously held in each chart's
    values.yaml `postgresql:` block:
      cluster_name   - Cluster .metadata.name; also the
                       prefix for the CNPG-managed Services (<cluster>-rw / -ro)
                       and the auto-generated "<cluster>-superuser" secret.
      image          - Postgres container image.
      instances      - Number of Postgres instances (replicas).
      database       - Application database created by initdb bootstrap.
      app_user       - Owner role for the application database. CNPG owns this
                       role's password (it reconciles it from the <cluster>-app
                       secret), so Vault must NOT rotate it -- the sibling
                       bootstrap Job creates "<app_user>_app" for that instead.
      storage_size   - PVC size.
      storage_class  - StorageClass name; empty string => operator/cluster default.
      resources      - Pod resource requests/limits block, passed through verbatim.
      monitoring     - Whether CNPG should emit a PodMonitor for Prometheus.
      cdc            - When true, provisions the Debezium logical-replication
                       role + uuid-ossp extension and sets wal_level=logical
                       (only catalog feeds the CDC -> OpenSearch pipeline).
      debezium_secret_name - Name of the basic-auth secret holding the debezium
                       role password (only consumed when cdc = true). Created by
                       the sibling `secret` module as "<svc>-postgres-debezium".
  EOT
  type = map(object({
    cluster_name         = string
    image                = string
    instances            = number
    database             = string
    app_user             = string
    storage_size         = string
    storage_class        = optional(string, "")
    resources            = any
    monitoring           = optional(bool, true)
    cdc                  = optional(bool, false)
    debezium_secret_name = optional(string, "")
  }))
}

variable "cross_service_readonly" {
  description = <<-EOT
    Cross-service read-only Postgres roles to provision on an existing (target_service)
    CNPG cluster. Each entry runs a bootstrap Job (as the CNPG superuser) that creates
    the LOGIN role and grants CONNECT, USAGE ON SCHEMA public and SELECT only on the
    named tables of that cluster's application database. Tables that do not exist yet
    are skipped with a warning; editing `tables` re-runs the Job.
      name           - reader/consumer service
      username       - Postgres role to create
      target_service - key in clusters of the cluster to read
      tables         - tables (schema public) to grant SELECT on
  EOT
  type = list(object({
    name           = string
    username       = string
    target_service = string
    tables         = list(string)
  }))
  default = []
}
