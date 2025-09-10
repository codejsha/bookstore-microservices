output "hosted_name" {
  description = "Created hosted repository name — used to derive its content privileges."
  value       = nexus_repository_npm_hosted.npm_hosted.name
}

output "proxy_name" {
  description = "Created proxy repository name."
  value       = nexus_repository_npm_proxy.npm_proxy.name
}

output "group_name" {
  description = "Created group repository name."
  value       = nexus_repository_npm_group.npm_group.name
}

output "maven_hosted_name" {
  description = "Created maven hosted repository name — publish target for the private platform artifacts."
  value       = nexus_repository_maven_hosted.maven_hosted.name
}

output "maven_proxy_name" {
  description = "Created maven proxy repository name."
  value       = nexus_repository_maven_proxy.maven_proxy.name
}

output "maven_group_name" {
  description = "Created maven group repository name — what the Gradle builds resolve against."
  value       = nexus_repository_maven_group.maven_group.name
}

output "pypi_group_name" {
  description = "Created pypi group repository name — what uv/pip resolve against."
  value       = nexus_repository_pypi_group.pypi_group.name
}

output "raw_hosted_name" {
  description = "Created raw hosted repository name — publish target for the vendored API specs."
  value       = nexus_repository_raw_hosted.raw_hosted.name
}

output "go_group_name" {
  description = "Created go group repository name — what GOPROXY resolves against."
  value       = nexus_repository_go_group.go_group.name
}

output "docker_proxy_name" {
  description = "Created docker proxy repository name — the pull-through cache for ghcr.io."
  value       = nexus_repository_docker_proxy.docker_proxy.name
}
