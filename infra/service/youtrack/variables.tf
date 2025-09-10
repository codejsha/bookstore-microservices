variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "chart_version" {
  description = "twenty20/helm-charts youtrack chart version"
  type        = string
}

variable "youtrack_version" {
  description = "YouTrack application image tag (jetbrains/youtrack)"
  type        = string
}

variable "youtrack_address" {
  description = "YouTrack public address (Ingress host + cert SAN)"
  type        = string
}

variable "storage_class" {
  description = "Default storage class for YouTrack PVCs (data/logs/conf/backup)"
  type        = string
}

variable "data_storage_size" {
  description = "Size of the YouTrack data PVC"
  type        = string
}

variable "logs_storage_size" {
  description = "Size of the YouTrack logs PVC"
  type        = string
}

variable "backup_storage_size" {
  description = "Size of the YouTrack backup PVC (volumeStorage mode)"
  type        = string
}

variable "admin_email" {
  description = "Operator email (used in Vault KV metadata)"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}
