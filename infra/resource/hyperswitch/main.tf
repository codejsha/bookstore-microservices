terraform {
  required_providers {
    restapi = {
      source  = "Mastercard/restapi"
      version = "~> 3.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

ephemeral "vault_kv_secret_v2" "app" {
  mount = "kv"
  name  = "hyperswitch/app/credentials"
}

provider "restapi" {
  uri                   = var.hyperswitch_admin_url
  write_returns_object  = true
  create_returns_object = true
  headers = {
    "api-key"      = ephemeral.vault_kv_secret_v2.app.data["admin_api_key"]
    "Content-Type" = "application/json"
  }
}

resource "restapi_object" "organization" {
  path         = "/organization"
  id_attribute = "organization_id"
  data = jsonencode({
    organization_name = var.organization_name
  })
  lifecycle {
    ignore_changes = all
  }
}

resource "restapi_object" "merchant" {
  path         = "/accounts"
  id_attribute = "merchant_id"
  data = jsonencode({
    merchant_id     = var.merchant_id
    merchant_name   = var.merchant_name
    organization_id = restapi_object.organization.id
  })
  lifecycle {
    ignore_changes = all
  }
}

resource "restapi_object" "api_key" {
  path         = "/api_keys/${restapi_object.merchant.id}"
  id_attribute = "key_id"
  data = jsonencode({
    name       = "payment-service"
    expiration = "never"
  })
  lifecycle {
    ignore_changes = all
  }
}

resource "random_password" "webhook_secret" {
  length  = 64
  special = false
}

resource "restapi_object" "merchant_webhook" {
  path           = "/accounts"
  object_id      = restapi_object.merchant.id
  id_attribute   = "merchant_id"
  create_path    = "/accounts/{id}"
  create_method  = "POST"
  read_path      = "/accounts/{id}"
  update_path    = "/accounts/{id}"
  update_method  = "POST"
  destroy_path   = "/accounts/{id}"
  destroy_method = "GET"

  ignore_all_server_changes = true

  data = jsonencode({
    merchant_id                  = restapi_object.merchant.id
    enable_payment_response_hash = true
    payment_response_hash_key    = random_password.webhook_secret.result
    webhook_details = {
      webhook_url              = var.payment_webhook_url
      payment_statuses_enabled = var.webhook_payment_statuses
      refund_statuses_enabled  = var.webhook_refund_statuses
    }
  })
}

resource "vault_kv_secret_v2" "payment" {
  mount = "kv"
  name  = "hyperswitch/payment"
  data_json = jsonencode({
    api_key        = jsondecode(restapi_object.api_key.create_response)["api_key"]
    key_id         = restapi_object.api_key.id
    profile_id     = jsondecode(restapi_object.merchant.create_response)["default_profile"]
    webhook_secret = random_password.webhook_secret.result
  })

  lifecycle {
    precondition {
      condition     = can(jsondecode(restapi_object.merchant.create_response)["default_profile"]) && jsondecode(restapi_object.merchant.create_response)["default_profile"] != ""
      error_message = "Hyperswitch merchant create_response has no non-empty 'default_profile'; refusing to write an empty profile_id to kv/hyperswitch/payment."
    }
  }
}
