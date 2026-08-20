variable "namespace" {
  description = "Bookstore application namespace"
  type        = string
}

variable "ingress_gateway_principal" {
  description = "SPIFFE principal of the north-south ingress gateway, allowed to reach directly ingress-facing services."
  type        = string
}

variable "hyperswitch_router_principal" {
  description = "SPIFFE principal of the Hyperswitch router, the only caller allowed onto payment's /internal/webhooks/* receiver."
  type        = string
}

variable "keycloak_issuer" {
  description = "Keycloak realm issuer URL used by Istio JWT validation"
  type        = string
}

variable "keycloak_jwks_uri" {
  description = "Keycloak JWKS endpoint used by Istio JWT validation"
  type        = string
}

variable "keycloak_audiences" {
  description = "Accepted audiences in the JWT 'aud' claim. Leave empty to skip audience validation."
  type        = list(string)
}

variable "anonymous_paths" {
  description = "Path patterns that bypass JWT enforcement (health/probes/internal-only endpoints)."
  type        = list(string)
}

variable "anonymous_read_paths" {
  description = "Path patterns anonymous clients may GET without a JWT (public catalog browsing); every other method on these paths stays enforced."
  type        = list(string)
  default     = []
}

variable "request_timeout" {
  description = "Mesh-wide per-request timeout enforced by the waypoint for every service."
  type        = string
}
