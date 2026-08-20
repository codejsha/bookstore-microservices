variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "chart_version" {
  description = "twenty20/helm-charts youtrack chart version"
  type        = string
}

variable "youtrack_version" {
  description = "YouTrack image tag"
  type        = string
}

variable "youtrack_address" {
  description = "YouTrack public address (Ingress host)"
  type        = string
}

variable "storage_class" {
  description = "StorageClass for PVCs"
  type        = string
}

variable "data_storage_size" {
  description = "Size of the data PVC"
  type        = string
}

variable "logs_storage_size" {
  description = "Size of the logs PVC"
  type        = string
}

variable "backup_storage_size" {
  description = "Size of the backup PVC (volumeStorage mode)"
  type        = string
}
