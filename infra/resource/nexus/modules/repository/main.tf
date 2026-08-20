terraform {
  required_providers {
    nexus = {
      source = "datadrivers/nexus"
    }
  }
}

resource "nexus_repository_npm_proxy" "npm_proxy" {
  name   = var.proxy_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = var.npm_registry_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_npm_hosted" "npm_hosted" {
  name   = var.hosted_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
    write_policy                   = var.hosted_write_policy
  }
}

resource "nexus_repository_npm_group" "npm_group" {
  name   = var.group_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  group {
    member_names = [
      nexus_repository_npm_hosted.npm_hosted.name,
      nexus_repository_npm_proxy.npm_proxy.name,
    ]
  }
}

resource "nexus_repository_maven_proxy" "maven_proxy" {
  name   = var.maven_proxy_name
  online = true

  maven {
    version_policy = "RELEASE"
    layout_policy  = "STRICT"
  }

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = var.maven_registry_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_maven_hosted" "maven_hosted" {
  name   = var.maven_hosted_name
  online = true

  maven {
    version_policy = "MIXED"
    layout_policy  = "STRICT"
  }

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
    write_policy                   = var.hosted_write_policy
  }
}

resource "nexus_repository_maven_proxy" "gradle_plugins_proxy" {
  name   = var.maven_plugins_proxy_name
  online = true

  maven {
    version_policy = "RELEASE"
    layout_policy  = "PERMISSIVE"
  }

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = false
  }

  proxy {
    remote_url       = var.maven_plugins_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_pypi_proxy" "pypi_proxy" {
  name   = var.pypi_proxy_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = var.pypi_registry_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_pypi_group" "pypi_group" {
  name   = var.pypi_group_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  group {
    member_names = [
      nexus_repository_pypi_proxy.pypi_proxy.name,
    ]
  }
}

resource "nexus_repository_maven_group" "maven_group" {
  name   = var.maven_group_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  group {
    member_names = [
      nexus_repository_maven_hosted.maven_hosted.name,
      nexus_repository_maven_proxy.maven_proxy.name,
      nexus_repository_maven_proxy.gradle_plugins_proxy.name,
    ]
  }
}

resource "nexus_repository_go_proxy" "go_proxy" {
  name   = var.go_proxy_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = var.go_registry_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_go_group" "go_group" {
  name   = var.go_group_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  group {
    member_names = [
      nexus_repository_go_proxy.go_proxy.name,
    ]
  }
}

resource "nexus_repository_docker_proxy" "docker_proxy" {
  name   = var.docker_proxy_name
  online = true

  docker {
    force_basic_auth = true
    v1_enabled       = false
    http_port        = var.docker_http_port
  }

  docker_proxy {
    index_type = "REGISTRY"
    index_url  = var.docker_registry_remote_url
  }

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = var.docker_registry_remote_url
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_raw_hosted" "raw_hosted" {
  name   = var.raw_hosted_name
  online = true

  storage {
    blob_store_name                = var.blob_store_name
    strict_content_type_validation = false
    write_policy                   = var.hosted_write_policy
  }
}
