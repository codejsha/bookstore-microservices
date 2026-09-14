terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

data "kubernetes_secret_v1" "cnpg_superuser" {
  for_each = toset(var.postgres_services)

  metadata {
    name      = "${each.key}-postgres-superuser"
    namespace = var.namespace
  }
}

resource "vault_database_secret_backend_connection" "postgres" {
  for_each = toset(var.postgres_services)

  backend = "database"
  name    = "${each.key}-postgres"

  password_policy = "password-alphanumeric"

  postgresql {
    connection_url = "postgresql://{{username}}:{{password}}@${each.key}-postgres-rw.${var.namespace}.svc.cluster.local:5432/${var.postgres_db_map[each.key]}?sslmode=disable"

    username = data.kubernetes_secret_v1.cnpg_superuser[each.key].data["username"]
    password = data.kubernetes_secret_v1.cnpg_superuser[each.key].data["password"]
  }

  allowed_roles = [
    "${each.key}-postgres-dynamic",
    "${each.key}-postgres-static",
  ]

  verify_connection = false
}

resource "vault_database_secret_backend_role" "postgres_dynamic" {
  for_each = toset(var.postgres_services)

  backend = "database"
  name    = "${each.key}-postgres-dynamic"
  db_name = vault_database_secret_backend_connection.postgres[each.key].name

  default_ttl = 3600
  max_ttl     = 86400

  creation_statements = [
    chomp(<<-SQL
      CREATE ROLE "{{name}}"
        WITH LOGIN
             PASSWORD '{{password}}'
             VALID UNTIL '{{expiration}}';
      GRANT USAGE ON SCHEMA public TO "{{name}}";
      GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO "{{name}}";
      ALTER DEFAULT PRIVILEGES IN SCHEMA public
        GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "{{name}}";
    SQL
    )
  ]

  revocation_statements = [
    chomp(<<-SQL
      REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
      REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{name}}";
      REVOKE USAGE ON SCHEMA public FROM "{{name}}";
      DROP ROLE IF EXISTS "{{name}}";
    SQL
    )
  ]
}

resource "vault_database_secret_backend_static_role" "postgres_static" {
  for_each = toset(var.postgres_services)

  backend  = "database"
  name     = "${each.key}-postgres-static"
  db_name  = vault_database_secret_backend_connection.postgres[each.key].name
  username = var.postgres_app_user_map[each.key]

  rotation_schedule = var.rotation_schedule
  rotation_window   = var.rotation_window
}

resource "terraform_data" "vault_agent_refresh" {
  triggers_replace = [var.rotation_schedule, var.rotation_window]

  depends_on = [vault_database_secret_backend_static_role.postgres_static]

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      NAMESPACE = var.namespace
      SERVICES  = join(" ", var.postgres_services)
    }
    command = <<-EOT
      set -u
      for svc in $SERVICES; do
        pods=$(kubectl -n "$NAMESPACE" get pods -l app.kubernetes.io/name="$svc" -o name 2>/dev/null | grep -v -- '-batch-' || true)
        if [ -z "$pods" ]; then
          echo "vault-agent-refresh: $svc has no pods, skipping"
          continue
        fi
        for p in $pods; do
          if kubectl -n "$NAMESPACE" exec "$p" -c vault-agent -- sh -c 'kill -TERM $(pgrep -x vault)' >/dev/null 2>&1; then
            echo "vault-agent-refresh: restarted vault-agent in $p"
          else
            echo "vault-agent-refresh: WARN could not restart vault-agent in $p" >&2
          fi
        done
      done
    EOT
  }
}
