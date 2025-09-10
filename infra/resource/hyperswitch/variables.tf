variable "vault_url" {
  type = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  type      = string
  sensitive = true
}

variable "hyperswitch_admin_url" {
  type        = string
  description = "Base URL of the Hyperswitch server admin API, reachable from the Terraform host."
}

variable "organization_name" {
  type    = string
  default = "Bookstore"
}

variable "merchant_id" {
  type        = string
  default     = "bookstore"
  description = "Stable merchant_id created on the Hyperswitch router."
}

variable "merchant_name" {
  type    = string
  default = "Bookstore"
}

variable "payment_webhook_url" {
  type        = string
  default     = "http://payment.bookstore.svc.cluster.local:8080/internal/webhooks/hyperswitch"
  description = "In-cluster URL the Hyperswitch router POSTs outgoing webhooks to (payment's HMAC-verified receiver)."
}

variable "webhook_payment_statuses" {
  type = list(string)
  default = [
    "succeeded",
    "failed",
    "cancelled",
    "cancelled_post_capture",
    "processing",
    "requires_capture",
    "partially_captured",
    "partially_captured_and_capturable",
    "expired",
  ]
  description = "Payment intent statuses that trigger an outgoing webhook."
}

variable "webhook_refund_statuses" {
  type        = list(string)
  default     = ["succeeded", "failed"]
  description = "Refund statuses that trigger an outgoing webhook."
}
