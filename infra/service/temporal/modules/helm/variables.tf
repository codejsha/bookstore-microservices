variable "namespace" {
  description = "Kubernetes namespace for Temporal"
  type        = string
}

variable "db_connect_addr" {
  description = "host:port of the MySQL read-write endpoint for Temporal's SQL stores"
  type        = string
}

variable "db_secret_name" {
  description = "Secret (key `password`) holding the MySQL root password"
  type        = string
}

variable "db_user" {
  description = "MySQL user Temporal connects as (root, so it can create databases/schemas)"
  type        = string
}

variable "authorization_enabled" {
  description = "Render server.config.authorization (JWT claim mapper + default authorizer) into the server config"
  type        = bool
  default     = false
}
