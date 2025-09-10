variable "group_name" {
  description = "Grafana alert rule group name (one evaluation group per service)"
  type        = string
}

variable "folder_title" {
  description = "Grafana folder title the rule group is created under"
  type        = string
}

variable "parent_folder_uid" {
  description = "UID of the parent folder the rule folder is nested under (null = top-level)"
  type        = string
  default     = null
}

variable "datasource_uid" {
  description = "Default datasource UID the rules query (see grafana datasources provisioning); a rule may override it"
  type        = string
  default     = "prometheus"
}

variable "datasource_type" {
  description = "Default datasource type of datasource_uid: prometheus or loki (LogQL metric queries evaluated as instant queries); a rule may override it"
  type        = string
  default     = "prometheus"
  validation {
    condition     = contains(["prometheus", "loki"], var.datasource_type)
    error_message = "datasource_type must be prometheus or loki."
  }
}

variable "interval_seconds" {
  description = "Rule group evaluation interval"
  type        = number
  default     = 60
}

variable "rules" {
  description = "Alert rules to provision into the group"
  type = list(object({
    name            = string
    expr            = string
    for             = string
    severity        = string
    summary         = string
    description     = string
    no_data_state   = optional(string, "OK")
    datasource_uid  = optional(string)
    datasource_type = optional(string)
  }))
}
