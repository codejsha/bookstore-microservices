variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "mysql_services" {
  description = "MySQL services"
  type = list(string)
}
