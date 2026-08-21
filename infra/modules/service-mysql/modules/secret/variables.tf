variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "mysql_services" {
  description = "MySQL services"
  type = list(string)
}

variable "mysql_service_accounts" {
  description = "MySQL service accounts"
  type = list(string)
}

variable "mysql_root_password" {
  description = "MySQL root password"
  type        = string
  sensitive   = true
}

variable "mysql_username" {
  description = "MySQL username"
  type        = string
  sensitive   = true
}

variable "mysql_password" {
  description = "MySQL password"
  type        = string
  sensitive   = true
}
