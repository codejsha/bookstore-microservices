variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "chart_version" {
  description = "Kyverno Helm chart version"
  type        = string
}

variable "image_pull_secret" {
  description = "dockerconfigjson Secret name passed to Kyverno as --imagePullSecrets"
  type        = string
}

variable "ca_bundle_configmap" {
  description = "CA bundle ConfigMap name distributed by trust-manager"
  type        = string
}

variable "ca_bundle_key" {
  description = "Key inside the CA bundle ConfigMap"
  type        = string
}
