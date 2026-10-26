variable "namespace" {
  description = "Kubernetes namespace for Temporal"
  type        = string
}

variable "db_connect_addr" {
  description = "host:port of the PostgreSQL read-write endpoint for Temporal's SQL stores"
  type        = string
}

variable "db_secret_name" {
  description = "Secret (key `password`) holding the PostgreSQL superuser password"
  type        = string
}

variable "db_user" {
  description = "PostgreSQL user Temporal connects as (superuser, so it can create databases/schemas)"
  type        = string
}

variable "authorization_enabled" {
  description = "Render server.config.authorization (JWT claim mapper + default authorizer) into the server config"
  type        = bool
  default     = false
}
