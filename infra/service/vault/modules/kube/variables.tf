variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kube_api_server_address" {
  description = "Kubernetes API server address"
  type        = string
}

variable "kube_ca_crt" {
  description = "Kubernetes CA certificate"
  type        = string
}

variable "tf_auth_roles" {
  description = "Map of role name => vault policy name. Creates a tf-<role> ServiceAccount + kubernetes-auth role per entry for out-of-cluster terraform runs."
  type        = map(string)
}
