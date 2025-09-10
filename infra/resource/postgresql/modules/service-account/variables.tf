variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "postgres_services" {
  description = "Postgres services"
  type        = list(string)
}
