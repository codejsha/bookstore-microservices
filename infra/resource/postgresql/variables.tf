variable "namespace" {
  description = "Namespace where the Postgres Cluster CRs live (catalog/customer/identity/inventory)"
  type        = string
}

variable "postgres_services" {
  description = "Services on Postgres (CloudNativePG managed)"
  type        = list(string)
}

variable "debezium_username" {
  description = "Debezium Postgres logical-replication role name (username only — password is generated)"
  type        = string
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

variable "postgres_db_map" {
  description = <<-EOT
    Map of service name → Postgres database name as set in each service's
    values.yaml (postgresql.database). Passed through to the db-engine module
    to build the Vault database connection URL.
    VERIFY AT APPLY: values must exactly match Cluster.spec.bootstrap.initdb.database.
  EOT
  type        = map(string)
}

variable "postgres_app_user_map" {
  description = <<-EOT
    Map of service name → the LOGIN role the Vault static-role rotates. This is
    NOT the CNPG initdb owner — it is the separate "<owner>_app" account created
    by the modules/cluster app-user Job (which does `CREATE ROLE <owner>_app`
    and `GRANT <owner> TO <owner>_app`). The static role logs in and rotates the
    password of this "_app" role.
    INVARIANT: each value MUST equal `<postgres_cluster_config[svc].app_user>_app`
    Do NOT set it to the bare owner name — that would point the static role at
    the wrong role.
  EOT
  type        = map(string)
}

variable "postgres_cluster_config" {
  description = <<-EOT
    Per-service CloudNativePG Cluster spec, keyed by service name. Replaces the
    `postgresql:` block that previously lived in each chart's values.yaml.
    INVARIANTS:
      - `database` MUST match postgres_db_map[svc] (both become initdb.database
        + the Vault db-engine connection URL).
      - `app_user` is the initdb.OWNER role. The app-user Job then derives the
        login account "<app_user>_app" from it; that "_app" account
        — NOT this owner — is what postgres_app_user_map points the Vault
        static role at. So app_user here does NOT equal postgres_app_user_map;
        rather postgres_app_user_map[svc] == "<app_user>_app".
    Set cdc = true only for the catalog cluster (Debezium logical replication).
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
