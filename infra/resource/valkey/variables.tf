variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "valkey_services" {
  description = "Valkey instances keyed by owning service"
  type = map(object({
    storage_size = string
    memory_limit = string
  }))
}

variable "chart_version" {
  description = "Valkey chart version"
  type        = string
}

variable "image_registry" {
  description = "Valkey image registry"
  type        = string
}

variable "image_repository" {
  description = "Valkey image repository"
  type        = string
}

variable "image_tag" {
  description = "Valkey image tag"
  type        = string
}
