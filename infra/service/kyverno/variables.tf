variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "chart_version" {
  description = "Kyverno Helm chart version"
  type        = string
}

variable "image_pull_secret" {
  description = "dockerconfigjson Secret in this namespace that Kyverno uses to pull signatures and attestations from the private registry. Provisioned by resource/kyverno, not here"
  type        = string
}

variable "ca_bundle_configmap" {
  description = "trust-manager Bundle ConfigMap distributed into every namespace; holds the public roots plus the internal PKI chain"
  type        = string
}

variable "ca_bundle_key" {
  description = "Key inside the CA bundle ConfigMap"
  type        = string
}
