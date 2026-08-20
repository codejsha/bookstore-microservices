variable "username" {
  description = "Permanent master-realm admin username"
  type        = string
}

variable "email" {
  description = "Permanent master-realm admin email"
  type        = string
}

variable "admin_role_id" {
  description = "ID of the master realm `admin` role"
  type        = string
}

variable "keycloak_namespace" {
  description = "Namespace running the Keycloak StatefulSet"
  type        = string
}
