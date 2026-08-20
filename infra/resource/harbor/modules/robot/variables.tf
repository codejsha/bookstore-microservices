variable "robot_name" {
  description = "Robot account name (Harbor prefixes the full name with robot$)"
  type        = string
  default     = "tekton-ci"
}

variable "projects" {
  description = "Projects the robot may push to and pull from"
  type        = list(string)
}
