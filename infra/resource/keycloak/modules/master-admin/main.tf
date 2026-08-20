terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
    vault = {
      source = "hashicorp/vault"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "random_password" "admin" {
  length  = 32
  special = false
}

resource "keycloak_user" "admin" {
  realm_id       = "master"
  username       = var.username
  email          = var.email
  email_verified = true
  enabled        = true

  initial_password {
    value     = random_password.admin.result
    temporary = false
  }
}

resource "keycloak_user_roles" "admin" {
  realm_id   = "master"
  user_id    = keycloak_user.admin.id
  role_ids   = [var.admin_role_id]
  exhaustive = false
}

resource "vault_kv_secret_v2" "admin" {
  mount = "kv"
  name  = "keycloak/admin/credentials"
  data_json = jsonencode({
    username = keycloak_user.admin.username
    password = random_password.admin.result
  })
}

resource "terraform_data" "temporary_admin" {
  triggers_replace = [keycloak_user.admin.id]

  provisioner "local-exec" {
    interpreter = ["/bin/sh", "-c"]
    environment = {
      KEYCLOAK_NAMESPACE = var.keycloak_namespace
      ADMIN_USERNAME     = keycloak_user.admin.username
      ADMIN_PASSWORD     = random_password.admin.result
    }
    command = <<-EOT
      set -eu
      kubectl exec -i -n "$KEYCLOAK_NAMESPACE" sts/keycloak -- sh -s "$ADMIN_USERNAME" "$ADMIN_PASSWORD" <<'SCRIPT'
      set -eu
      kcadm="/opt/keycloak/bin/kcadm.sh"
      cfg=/tmp/kcadm.config
      trap 'rm -f "$cfg"' EXIT
      "$kcadm" config credentials --config "$cfg" --server http://localhost:8080 --realm master --user "$1" --password "$2"
      ids=$("$kcadm" get users --config "$cfg" -r master -q 'q=is_temporary_admin:true' --fields id --format csv --noquotes)
      for id in $ids; do
        echo "deleting temporary admin $id"
        "$kcadm" delete "users/$id" --config "$cfg" -r master
      done
      SCRIPT
    EOT
  }

  depends_on = [keycloak_user_roles.admin, vault_kv_secret_v2.admin]
}
