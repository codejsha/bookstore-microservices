variable "robot_name" {
  description = "Robot account name (Harbor prefixes the full name with robot$)"
  type        = string
  default     = "tekton-ci"
}

variable "projects" {
  description = "Projects the robot may push to and pull from"
  type        = list(string)
}

variable "pull_robot_name" {
  description = "Name of the system robot account used for cluster-wide image pulls"
  type        = string
  default     = "cluster-pull"
}
