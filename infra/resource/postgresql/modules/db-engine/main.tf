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
    <<-SQL
      CREATE ROLE "{{name}}"
        WITH LOGIN
             PASSWORD '{{password}}'
             VALID UNTIL '{{expiration}}';
      GRANT USAGE ON SCHEMA public TO "{{name}}";
      GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO "{{name}}";
      ALTER DEFAULT PRIVILEGES IN SCHEMA public
        GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "{{name}}";
    SQL
  ]

  revocation_statements = [
    <<-SQL
      REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
      REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{name}}";
      REVOKE USAGE ON SCHEMA public FROM "{{name}}";
      DROP ROLE IF EXISTS "{{name}}";
    SQL
  ]
}

resource "vault_database_secret_backend_static_role" "postgres_static" {
  for_each = toset(var.postgres_services)

  backend  = "database"
  name     = "${each.key}-postgres-static"
  db_name  = vault_database_secret_backend_connection.postgres[each.key].name
  username = var.postgres_app_user_map[each.key]

  rotation_period = 86400
}
