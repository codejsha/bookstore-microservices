variable "proxy_name" {
  description = "npm proxy repository name"
  type        = string
}

variable "hosted_name" {
  description = "npm hosted repository name"
  type        = string
}

variable "group_name" {
  description = "npm group repository name"
  type        = string
}

variable "blob_store_name" {
  description = "Blob store backing the repositories"
  type        = string
}

variable "npm_registry_remote_url" {
  description = "Upstream npm registry the proxy caches"
  type        = string
}

variable "hosted_write_policy" {
  description = "Hosted repo deployment policy (ALLOW | ALLOW_ONCE | DENY)"
  type        = string
}

variable "maven_proxy_name" {
  description = "maven proxy repository name"
  type        = string
}

variable "maven_hosted_name" {
  description = "maven hosted repository name"
  type        = string
}

variable "maven_group_name" {
  description = "maven group repository name"
  type        = string
}

variable "maven_registry_remote_url" {
  description = "Upstream maven registry the proxy caches"
  type        = string
}

variable "maven_plugins_proxy_name" {
  description = "maven proxy repository name for the Gradle Plugin Portal"
  type        = string
}

variable "maven_plugins_remote_url" {
  description = "Gradle Plugin Portal m2 URL the proxy caches"
  type        = string
}

variable "pypi_proxy_name" {
  description = "pypi proxy repository name (caches pypi.org)"
  type        = string
}

variable "pypi_group_name" {
  description = "pypi group repository name — what uv/pip resolve against"
  type        = string
}

variable "pypi_registry_remote_url" {
  description = "Upstream PyPI registry the proxy caches"
  type        = string
}

variable "raw_hosted_name" {
  description = "Raw hosted repository holding the vendored third-party API specs."
  type        = string
}

variable "go_proxy_name" {
  description = "go proxy repository name (caches proxy.golang.org)"
  type        = string
}

variable "go_group_name" {
  description = "go group repository name — the single URL GOPROXY points at"
  type        = string
}

variable "go_registry_remote_url" {
  description = "Upstream Go module proxy the proxy caches"
  type        = string
}

variable "docker_proxy_name" {
  description = "docker proxy repository name (caches the ghcr.io OCI artifacts CI pulls)"
  type        = string
}

variable "docker_registry_remote_url" {
  description = "Upstream OCI registry the docker proxy caches"
  type        = string
}

variable "docker_http_port" {
  description = <<-EOT
    Dedicated HTTP connector port for the docker proxy. Nexus OSS has no
    path-based docker routing, so a docker repo is only reachable on its own
    port; infra/service/nexus fronts that connector with the nexus-docker Service.
  EOT
  type        = number
}
