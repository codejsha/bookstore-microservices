variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "istio_version" {
  description = "Version for all four Istio charts (base, istiod, cni, ztunnel) -- they must move in lockstep. Keep kiali_chart_version compatible when bumping."
  type        = string
}
