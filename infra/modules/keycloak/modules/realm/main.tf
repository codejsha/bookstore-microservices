terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "keycloak_realm" {
  name = "keycloak-realm"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "keycloak_realm" {
  role_name = "keycloak-realm-role"
  bound_service_account_names = ["keycloak-sa"]
  bound_service_account_namespaces = [var.namespace]
  token_policies = [vault_policy.keycloak_realm.name]
}

resource "vault_kv_secret_v2" "keycloak_realm" {
  depends_on = [vault_kubernetes_auth_backend_role.keycloak_realm]
  name  = "keycloak/realm"
  mount = "kv"
  data_json = file("${path.module}/realm-export.json")
}
