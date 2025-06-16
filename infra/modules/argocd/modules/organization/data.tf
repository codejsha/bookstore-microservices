data "vault_generic_secret" "catalog_helm_ssh_keys" {
  path = "kv/gitea/ssh/catalog-helm"
}

data "vault_generic_secret" "customer_helm_ssh_keys" {
  path = "kv/gitea/ssh/customer-helm"
}

data "vault_generic_secret" "identity_helm_ssh_keys" {
  path = "kv/gitea/ssh/identity-helm"
}

data "vault_generic_secret" "inventory_helm_ssh_keys" {
  path = "kv/gitea/ssh/inventory-helm"
}

data "vault_generic_secret" "order_helm_ssh_keys" {
  path = "kv/gitea/ssh/order-helm"
}

data "vault_generic_secret" "payment_helm_ssh_keys" {
  path = "kv/gitea/ssh/payment-helm"
}
