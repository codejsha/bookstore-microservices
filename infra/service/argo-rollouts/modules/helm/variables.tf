variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gatewayapi_plugin_url" {
  description = "Download location of the Gateway API trafficRouter plugin binary (linux/amd64). The controller fetches it into the plugin-bin emptyDir at startup."
  type        = string
}
