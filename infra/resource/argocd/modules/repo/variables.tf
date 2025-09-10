variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "gitea_ssh_fqdn" {
  description = "Gitea SSH service FQDN. Must match the host in the argocd-ssh-known-hosts-cm entry, otherwise the host key check fails."
  type        = string
}

variable "gitea_ssh_port" {
  description = "Gitea SSH service port"
  type        = number
  default     = 22
}

variable "app_repos" {
  description = "List of repositories"
  type        = list(string)
}
