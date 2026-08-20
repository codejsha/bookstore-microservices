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

data "kubernetes_secret_v1" "mysql_root" {
  for_each = toset(var.mysql_services)
  metadata {
    name      = var.mysql_db_config[each.key].secret_name
    namespace = var.namespace
  }
}

locals {
  readonly_entries = {
    for e in var.cross_service_readonly : "${e.name}-${e.target_service}" => e
  }

  readonly_static_roles_by_target = {
    for tgt in distinct([for e in var.cross_service_readonly : e.target_service]) :
    tgt => [
      for e in var.cross_service_readonly :
      "${e.name}-${e.target_service}-mysql-static" if e.target_service == tgt
    ]
  }
}

resource "vault_database_secret_backend_connection" "mysql" {
  for_each = toset(var.mysql_services)

  backend     = "database"
  name        = "${each.key}-mysql"
  plugin_name = "mysql-database-plugin"

  password_policy = "password-alphanumeric"

  allowed_roles = concat(
    ["${each.key}-mysql-dynamic", "${each.key}-mysql-static"],
    lookup(local.readonly_static_roles_by_target, each.key, []),
  )

  verify_connection = false

  mysql {
    connection_url = "{{username}}:{{password}}@tcp(${var.mysql_db_config[each.key].cluster_name}.${var.namespace}.svc.cluster.local:3306)/"

    username = data.kubernetes_secret_v1.mysql_root[each.key].data["rootUser"]
    password = data.kubernetes_secret_v1.mysql_root[each.key].data["rootPassword"]
  }
}

resource "vault_database_secret_backend_role" "mysql_dynamic" {
  for_each = toset(var.mysql_services)

  backend = "database"
  name    = "${each.key}-mysql-dynamic"

  db_name = vault_database_secret_backend_connection.mysql[each.key].name

  default_ttl = 3600
  max_ttl     = 86400

  creation_statements = [
    "CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';",
    "GRANT ALL PRIVILEGES ON `${var.mysql_db_config[each.key].database_name}`.* TO '{{name}}'@'%';",
    "FLUSH PRIVILEGES;",
  ]

  revocation_statements = [
    "DROP USER IF EXISTS '{{name}}'@'%';",
  ]

  rollback_statements = [
    "DROP USER IF EXISTS '{{name}}'@'%';",
  ]
}

resource "vault_database_secret_backend_static_role" "mysql_static" {
  for_each = toset(var.mysql_services)

  backend  = "database"
  name     = "${each.key}-mysql-static"
  db_name  = vault_database_secret_backend_connection.mysql[each.key].name
  username = var.mysql_db_config[each.key].app_user

  rotation_period = 86400
  rotation_statements = [
    "ALTER USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';",
  ]
}

resource "vault_database_secret_backend_static_role" "cross_service_readonly" {
  for_each = local.readonly_entries

  backend  = "database"
  name     = "${each.key}-mysql-static"
  db_name  = vault_database_secret_backend_connection.mysql[each.value.target_service].name
  username = each.value.username

  rotation_period = 86400
  rotation_statements = [
    "ALTER USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';",
  ]
}
