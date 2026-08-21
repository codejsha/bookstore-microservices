variable "namespace" {
  description = "Namespace name"
  type        = string
  sensitive   = true
}

variable "host_address" {
  description = "Host address"
  type        = string
}

variable "host_fqdn" {
  description = "Host FQDN"
  type        = string
}

variable "dest_port" {
  description = "Destination port"
  type        = number
}

variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}
