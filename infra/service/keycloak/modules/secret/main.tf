terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    vault = {
      source = "hashicorp/vault"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

resource "null_resource" "wait_for_initial_admin" {
  provisioner "local-exec" {
    command = <<EOT
      for i in $(seq 1 90); do
        kubectl get secret keycloak-initial-admin -n ${var.namespace} >/dev/null 2>&1 && exit 0
        sleep 10
      done
      echo "timed out waiting for secret keycloak-initial-admin in ${var.namespace}" >&2
      exit 1
    EOT
  }
}

data "kubernetes_secret_v1" "keycloak_initial_admin" {
  depends_on = [null_resource.wait_for_initial_admin]
  metadata {
    name      = "keycloak-initial-admin"
    namespace = var.namespace
  }
}

removed {
  from = vault_kv_secret_v2.keycloak_admin

  lifecycle {
    destroy = false
  }
}

resource "vault_kv_secret_v2" "keycloak_bootstrap_admin" {
  mount = "kv"
  name  = "keycloak/bootstrap/credentials"
  data_json = jsonencode({
    username = data.kubernetes_secret_v1.keycloak_initial_admin.data["username"]
    password = data.kubernetes_secret_v1.keycloak_initial_admin.data["password"]
  })
}
