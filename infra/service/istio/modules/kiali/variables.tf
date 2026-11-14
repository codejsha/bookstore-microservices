variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kiali_chart_version" {
  description = "kiali-server helm chart version. Kiali tests each release against the currently-supported Istio releases; 2.26.0 covers Istio 1.29. See https://kiali.io/docs/installation/installation-guide/prerequisites/"
  type        = string
}

variable "grafana_username" {
  description = "Grafana username"
  type        = string
  sensitive   = true
}

variable "grafana_password" {
  description = "Grafana password"
  type        = string
  sensitive   = true
}

variable "oidc_group_prefix" {
  description = "Prefix the API server AuthenticationConfiguration adds to OIDC groups claim values"
  type        = string
  default     = "oidc:"
}

variable "viewer_roles" {
  description = "Platform realm roles granted read-only mesh access through Kiali"
  type        = list(string)
  default     = ["PLATFORM_DEV", "PLATFORM_DEVOPS", "PLATFORM_ADMIN"]
}

variable "editor_roles" {
  description = "Platform realm roles granted cluster-wide mesh write access and pod port-forwarding through Kiali"
  type        = list(string)
  default     = ["PLATFORM_DEVOPS", "PLATFORM_ADMIN"]
}

variable "project_namespaces" {
  description = "Namespaces the bookstore project roles get mesh access to through Kiali"
  type        = list(string)
  default     = ["bookstore"]
}

variable "project_viewer_roles" {
  description = "Platform realm roles granted read-only mesh access to the project namespaces through Kiali"
  type        = list(string)
  default     = ["BOOKSTORE_DEV", "BOOKSTORE_DEVOPS", "BOOKSTORE_ADMIN"]
}

variable "project_editor_roles" {
  description = "Platform realm roles granted mesh write access and pod port-forwarding in the project namespaces through Kiali"
  type        = list(string)
  default     = ["BOOKSTORE_DEVOPS", "BOOKSTORE_ADMIN"]
}
