variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "gitea_fqdn" {
  description = "Gitea FQDN"
  type        = string
}

variable "app_repos" {
  description = "List of repositories"
  type = list(string)
}
