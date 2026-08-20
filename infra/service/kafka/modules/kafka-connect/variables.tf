variable "namespace" {
  description = "Namespace where KafkaConnect runs (same as Kafka broker)"
  type        = string
}

variable "kafka_cluster_name" {
  description = "Strimzi Kafka cluster name"
  type        = string
}

variable "connect_cluster_name" {
  description = "KafkaConnect cluster name"
  type        = string
}

variable "connect_image" {
  description = "Fully-qualified image ref Strimzi will push the built Connect image to"
  type        = string
}

variable "replicas" {
  description = "KafkaConnect worker replicas"
  type        = number
}

variable "image_pull_secret" {
  description = "imagePullSecret name (empty to skip)"
  type        = string
}

variable "additional_build_options" {
  description = "Extra container-build options. Must be in Strimzi's allow-list."
  type        = list(string)
}

variable "build_host_aliases" {
  description = "hostAliases for the Kaniko build pod (scoped /etc/hosts)."
  type = list(object({
    ip        = string
    hostnames = list(string)
  }))
}

variable "external_secret_mounts" {
  description = "K8s Secrets (same namespace) to mount into Connect pods for $${file:...} config provider refs"
  type = list(object({
    name       = string
    mount_name = string
  }))
}
