data "vault_generic_secret" "admin_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.admin_service}-helm"
}

data "vault_generic_secret" "catalog_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.catalog_service}-helm"
}

data "vault_generic_secret" "customer_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.customer_service}-helm"
}

data "vault_generic_secret" "identity_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.identity_service}-helm"
}

data "vault_generic_secret" "inventory_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.inventory_service}-helm"
}

data "vault_generic_secret" "order_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.order_service}-helm"
}

data "vault_generic_secret" "payment_helm_ssh_keys" {
  path = "kv/gitea/ssh/${local.payment_service}-helm"
}
