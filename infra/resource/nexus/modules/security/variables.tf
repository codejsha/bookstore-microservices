variable "role_id" {
  description = "Security role id for the publish identity"
  type        = string
}

variable "publisher_username" {
  description = "Publish user id"
  type        = string
}

variable "publisher_password" {
  description = "Publish user password (write-only on the Nexus user; persisted to Vault by the caller)"
  type        = string
  sensitive   = true
}

variable "hosted_repository_name" {
  description = "Hosted npm repository whose content privileges the role grants"
  type        = string
}

variable "group_repository_name" {
  description = "Group npm repository whose read privileges the role grants (so the publish identity can also `npm install` once anonymous access is disabled)"
  type        = string
}

variable "maven_hosted_repository_name" {
  description = "Hosted maven repository whose content privileges the role grants (publish target for the private platform artifacts)"
  type        = string
}

variable "maven_group_repository_name" {
  description = "Group maven repository whose read privileges the role grants (how the Gradle CI builds resolve dependencies, since anonymous access is disabled)"
  type        = string
}

variable "pypi_group_repository_name" {
  description = "Group pypi repository whose read privileges the role grants (uv/pip resolve through it; anonymous access is disabled)"
  type        = string
}

variable "raw_hosted_repository_name" {
  description = "Raw hosted repository holding the vendored API specs."
  type        = string
}

variable "go_group_repository_name" {
  description = "Group go repository whose read privileges the role grants (GOPROXY resolves through it; anonymous access is disabled)"
  type        = string
}

variable "docker_proxy_repository_name" {
  description = "Docker proxy repository whose read privileges the role grants (CI pulls the trivy DBs through it; anonymous access is disabled)"
  type        = string
}
