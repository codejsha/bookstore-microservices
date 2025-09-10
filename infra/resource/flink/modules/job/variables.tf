variable "namespace" {
  description = "Namespace for FlinkSessionJob"
  type        = string
}

variable "job_name" {
  description = "FlinkSessionJob name"
  type        = string
}

variable "session_cluster_name" {
  description = "FlinkDeployment (session cluster) name to submit to"
  type        = string
}

variable "job_jar_uri" {
  description = "JAR URI for the Flink job"
  type        = string
}

variable "job_parallelism" {
  description = "Job parallelism"
  type        = number
}

variable "job_args" {
  description = "Job arguments"
  type        = list(string)
}
