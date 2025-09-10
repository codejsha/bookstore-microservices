terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  cdc_post_init_sql = {
    for svc, c in var.clusters : svc => [
      for stmt in split("\n", templatefile("${path.module}/sql/cdc-post-init.sql.tpl", {
        database = c.database
        app_user = c.app_user
      })) : stmt if trimspace(stmt) != ""
    ] if c.cdc
  }

  cdc_postgresql_parameters = {
    wal_level             = "logical"
    max_wal_senders       = "10"
    max_replication_slots = "10"
    shared_buffers        = "256MB"
  }

  manifests = {
    for svc, c in var.clusters : svc => {
      apiVersion = "postgresql.cnpg.io/v1"
      kind       = "Cluster"
      metadata = {
        name      = c.cluster_name
        namespace = var.namespace
      }
      spec = merge(concat(
        [{
          instances             = c.instances
          imageName             = c.image
          primaryUpdateStrategy = "unsupervised"

          enableSuperuserAccess = true

          storage = merge(
            { size = c.storage_size },
            c.storage_class != "" ? { storageClass = c.storage_class } : {},
          )

          resources = c.resources

          bootstrap = {
            initdb = merge(
              {
                database = c.database
                owner    = c.app_user
              },
              c.cdc ? { postInitApplicationSQL = local.cdc_post_init_sql[svc] } : {},
            )
          }

          monitoring = {
            enablePodMonitor = c.monitoring
          }

          inheritedMetadata = {
            labels = {
              release = "prometheus"
            }
          }
        }],
        c.cdc ? [{
          postgresql = {
            parameters = local.cdc_postgresql_parameters
          }
          managed = {
            roles = [
              {
                name        = "debezium"
                ensure      = "present"
                login       = true
                replication = true
                inRoles     = [c.app_user]
                passwordSecret = {
                  name = c.debezium_secret_name
                }
              }
            ]
          }
        }] : [],
      )...)
    }
  }
}

resource "kubernetes_manifest" "cluster" {
  for_each = local.manifests
  manifest = each.value

  computed_fields = [
    "metadata.labels",
    "metadata.annotations",
    "spec.postgresql.parameters",
  ]
}

resource "kubernetes_job_v1" "app_user" {
  for_each = var.clusters

  depends_on = [kubernetes_manifest.cluster]

  metadata {
    name      = "${each.key}-postgres-app-user"
    namespace = var.namespace
  }

  spec {
    backoff_limit = 6
    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "${each.key}-postgres-app-user"
          "app.kubernetes.io/component" = "bootstrap"
        }
      }
      spec {
        restart_policy = "OnFailure"
        container {
          name  = "bootstrap"
          image = each.value.image

          env {
            name = "PGUSER"
            value_from {
              secret_key_ref {
                name = "${each.value.cluster_name}-superuser"
                key  = "username"
              }
            }
          }
          env {
            name = "PGPASSWORD"
            value_from {
              secret_key_ref {
                name = "${each.value.cluster_name}-superuser"
                key  = "password"
              }
            }
          }
          env {
            name  = "PGHOST"
            value = "${each.value.cluster_name}-rw.${var.namespace}.svc.cluster.local"
          }
          env {
            name  = "PGDATABASE"
            value = each.value.database
          }
          env {
            name  = "APP_OWNER"
            value = each.value.app_user
          }
          env {
            name  = "APP_USER"
            value = "${each.value.app_user}_app"
          }

          command = ["/bin/sh", "-c"]
          args = [<<-EOT
            set -eu
            echo "Waiting for Postgres at $PGHOST ..."
            i=0
            until pg_isready -h "$PGHOST" -U "$PGUSER" -q; do
              i=$((i+1))
              if [ "$i" -ge 60 ]; then
                echo "Postgres did not become ready in time" >&2
                exit 1
              fi
              echo "  not ready yet (attempt $i), retrying in 10s..."
              sleep 10
            done
            echo "Postgres is ready. Bootstrapping $APP_USER ..."

            psql -h "$PGHOST" -d "$PGDATABASE" -v ON_ERROR_STOP=1 <<SQL
            DO \$\$
            BEGIN
              IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$APP_USER') THEN
                EXECUTE format('CREATE ROLE %I WITH LOGIN PASSWORD %L', '$APP_USER', gen_random_uuid()::text);
              END IF;
            END
            \$\$;
            GRANT "$APP_OWNER" TO "$APP_USER";
            GRANT CONNECT ON DATABASE "$PGDATABASE" TO "$APP_USER";
            SQL
            echo "Bootstrap complete: $APP_USER inherits $APP_OWNER."
          EOT
          ]
        }
      }
    }
  }

  wait_for_completion = true
  timeouts {
    create = "10m"
    update = "10m"
  }
}
