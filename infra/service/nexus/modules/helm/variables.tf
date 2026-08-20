variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "docker_port" {
  description = <<-EOT
    Docker connector port. Must match the port in values.yaml's
    nexus.docker.registries — the chart publishes that as <fullname>-docker-<port>,
    a name callers cannot spell without knowing the port, so this module fronts
    the same pods with a stable nexus-docker Service.
  EOT
  type        = number
  default     = 5000
}
