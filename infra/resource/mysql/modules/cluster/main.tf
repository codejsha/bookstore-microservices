terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  readonly_users = {
    for e in var.cross_service_readonly : "${e.name}-${e.target_service}" => e
  }
}

resource "kubernetes_manifest" "innodb_cluster" {
  for_each = toset(var.mysql_services)

  manifest = {
    apiVersion = "mysql.oracle.com/v2"
    kind       = "InnoDBCluster"
    metadata = {
      name      = var.mysql_db_config[each.key].cluster_name
      namespace = var.namespace
    }
    spec = {
      instances        = var.mysql_db_config[each.key].instances
      secretName       = var.mysql_db_config[each.key].secret_name
      tlsUseSelfSigned = true
      version          = var.mysql_db_config[each.key].version
      router = {
        instances = 1
      }
      metrics = {
        enable  = true
        image   = var.mysqld_exporter_image
        monitor = true
      }
      datadirVolumeClaimTemplate = {
        accessModes = ["ReadWriteOnce"]
        resources = {
          requests = {
            storage = var.mysql_db_config[each.key].storage_size
          }
        }
      }
      podSpec = {
        containers = [
          {
            name      = "mysql"
            resources = var.mysql_resources
          }
        ]
      }
    }
  }
}

resource "kubernetes_job_v1" "bootstrap" {
  for_each = toset(var.mysql_services)

  wait_for_completion = true

  timeouts {
    create = "10m"
    update = "10m"
  }

  metadata {
    name      = "${each.key}-mysql-bootstrap"
    namespace = var.namespace
  }

  spec {
    backoff_limit           = 6
    active_deadline_seconds = 1800

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "${each.key}-mysql-bootstrap"
          "app.kubernetes.io/component" = "db-bootstrap"
        }
      }

      spec {
        restart_policy = "OnFailure"

        container {
          name  = "bootstrap"
          image = var.mysql_client_image

          env {
            name = "MYSQL_ROOT_USER"
            value_from {
              secret_key_ref {
                name = var.mysql_db_config[each.key].secret_name
                key  = "rootUser"
              }
            }
          }
          env {
            name = "MYSQL_ROOT_PASSWORD"
            value_from {
              secret_key_ref {
                name = var.mysql_db_config[each.key].secret_name
                key  = "rootPassword"
              }
            }
          }
          env {
            name  = "MYSQL_HOST"
            value = "${var.mysql_db_config[each.key].cluster_name}.${var.namespace}.svc.cluster.local"
          }
          env {
            name  = "MYSQL_PORT"
            value = "3306"
          }
          env {
            name  = "APP_DATABASE"
            value = var.mysql_db_config[each.key].database_name
          }
          env {
            name  = "APP_USER"
            value = var.mysql_db_config[each.key].app_user
          }

          command = ["/bin/sh", "-c"]
          args = [<<-EOT
            set -eu
            echo "Waiting for MySQL at $MYSQL_HOST:$MYSQL_PORT ..."
            i=0
            until mysqladmin ping \
                  -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
                  -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" \
                  --silent; do
              i=$((i+1))
              if [ "$i" -ge 60 ]; then
                echo "MySQL did not become ready in time" >&2
                exit 1
              fi
              echo "  not ready yet (attempt $i), retrying in 10s..."
              sleep 10
            done
            echo "MySQL is ready. Bootstrapping database and app user..."
            mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
                  -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" <<SQL
            CREATE DATABASE IF NOT EXISTS \`$APP_DATABASE\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
            CREATE USER IF NOT EXISTS '$APP_USER'@'%' IDENTIFIED BY RANDOM PASSWORD;
            GRANT ALL PRIVILEGES ON \`$APP_DATABASE\`.* TO '$APP_USER'@'%';
            FLUSH PRIVILEGES;
            SQL
            echo "Bootstrap complete for $APP_DATABASE / $APP_USER."
          EOT
          ]
        }
      }
    }
  }

  depends_on = [kubernetes_manifest.innodb_cluster]
}

resource "kubernetes_job_v1" "readonly_bootstrap" {
  for_each = local.readonly_users

  wait_for_completion = true

  timeouts {
    create = "10m"
    update = "10m"
  }

  metadata {
    name      = "${each.key}-ro-bootstrap"
    namespace = var.namespace
  }

  spec {
    backoff_limit           = 6
    active_deadline_seconds = 1800

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "${each.key}-ro-bootstrap"
          "app.kubernetes.io/component" = "db-readonly-bootstrap"
        }
      }

      spec {
        restart_policy = "OnFailure"

        container {
          name  = "bootstrap"
          image = var.mysql_client_image

          env {
            name = "MYSQL_ROOT_USER"
            value_from {
              secret_key_ref {
                name = var.mysql_db_config[each.value.target_service].secret_name
                key  = "rootUser"
              }
            }
          }
          env {
            name = "MYSQL_ROOT_PASSWORD"
            value_from {
              secret_key_ref {
                name = var.mysql_db_config[each.value.target_service].secret_name
                key  = "rootPassword"
              }
            }
          }
          env {
            name  = "MYSQL_HOST"
            value = "${var.mysql_db_config[each.value.target_service].cluster_name}.${var.namespace}.svc.cluster.local"
          }
          env {
            name  = "MYSQL_PORT"
            value = "3306"
          }
          env {
            name  = "APP_DATABASE"
            value = var.mysql_db_config[each.value.target_service].database_name
          }
          env {
            name  = "RO_USER"
            value = each.value.username
          }
          env {
            name  = "RO_TABLES"
            value = join(" ", each.value.tables)
          }

          command = ["/bin/sh", "-c"]
          args = [<<-EOT
            set -eu
            echo "Waiting for MySQL at $MYSQL_HOST:$MYSQL_PORT ..."
            i=0
            until mysqladmin ping \
                  -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
                  -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" \
                  --silent; do
              i=$((i+1))
              if [ "$i" -ge 60 ]; then
                echo "MySQL did not become ready in time" >&2
                exit 1
              fi
              echo "  not ready yet (attempt $i), retrying in 10s..."
              sleep 10
            done
            echo "MySQL is ready. Provisioning read-only user $RO_USER on $APP_DATABASE ..."
            {
              echo "CREATE USER IF NOT EXISTS '$RO_USER'@'%' IDENTIFIED BY RANDOM PASSWORD;"
              for t in $RO_TABLES; do
                echo "GRANT SELECT ON \`$APP_DATABASE\`.\`$t\` TO '$RO_USER'@'%';"
              done
              echo "FLUSH PRIVILEGES;"
            } | mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
                      -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD"
            echo "Read-only user $RO_USER provisioned with SELECT on: $RO_TABLES"
          EOT
          ]
        }
      }
    }
  }

  depends_on = [kubernetes_job_v1.bootstrap]
}
