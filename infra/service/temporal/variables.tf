variable "namespace" {
  description = "Kubernetes namespace for Temporal"
  type        = string
}

variable "temporal_address" {
  description = "Temporal Web UI public address (HTTPRoute hostname)"
  type        = string
}

variable "vault_url" {
  description = "Vault address for the kubernetes-auth provider login"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "SA token for Vault kubernetes-auth (kubectl create token tf-temporal)"
  type        = string
  sensitive   = true
}

variable "mysql_cluster_name" {
  description = "InnoDBCluster name (also the RW Service name on :3306)"
  type        = string
}

variable "mysql_instances" {
  description = "MySQL server instances (3 for Group Replication quorum; 1 = no HA, cannot self-heal after a restart)"
  type        = number
}

variable "mysql_router_instances" {
  description = "MySQL Router instances"
  type        = number
}

variable "mysql_version" {
  description = "MySQL server version"
  type        = string
}

variable "mysql_storage_size" {
  description = "Per-instance data volume size"
  type        = string
}

variable "mysql_storage_class_name" {
  description = "StorageClass for the data volume (k3s default is local-path)"
  type        = string
}

variable "mysql_db_secret_name" {
  description = "Secret holding the password under key `password` for the Temporal chart's existingSecret"
  type        = string
}

variable "mysql_db_user" {
  description = "MySQL user Temporal connects as (root, so it can create databases/schemas)"
  type        = string
}

variable "mysql_resources" {
  description = "Resource requests/limits for the MySQL server container"
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
}

variable "mysql_router_resources" {
  description = "Resource requests/limits for the MySQL Router container"
  type = object({
    requests = object({ cpu = string, memory = string })
    limits   = object({ cpu = string, memory = string })
  })
}
