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

resource "terraform_data" "readonly_tables" {
  for_each = local.readonly_users

  input = join(" ", each.value.tables)
}

resource "kubernetes_job_v1" "readonly_bootstrap" {
  for_each = local.readonly_users

  depends_on = [kubernetes_job_v1.app_user]

  lifecycle {
    replace_triggered_by = [terraform_data.readonly_tables[each.key]]
  }

  metadata {
    name      = "${each.key}-postgres-readonly"
    namespace = var.namespace
  }

  spec {
    backoff_limit = 6
    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "${each.key}-postgres-readonly"
          "app.kubernetes.io/component" = "readonly-bootstrap"
        }
      }
      spec {
        restart_policy = "OnFailure"
        container {
          name  = "bootstrap"
          image = var.clusters[each.value.target_service].image

          env {
            name = "PGUSER"
            value_from {
              secret_key_ref {
                name = "${var.clusters[each.value.target_service].cluster_name}-superuser"
                key  = "username"
              }
            }
          }
          env {
            name = "PGPASSWORD"
            value_from {
              secret_key_ref {
                name = "${var.clusters[each.value.target_service].cluster_name}-superuser"
                key  = "password"
              }
            }
          }
          env {
            name  = "PGHOST"
            value = "${var.clusters[each.value.target_service].cluster_name}-rw.${var.namespace}.svc.cluster.local"
          }
          env {
            name  = "PGDATABASE"
            value = var.clusters[each.value.target_service].database
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
            echo "Postgres is ready. Provisioning read-only role $RO_USER on $PGDATABASE ..."

            psql -h "$PGHOST" -d "$PGDATABASE" -v ON_ERROR_STOP=1 <<SQL
            DO \$\$
            BEGIN
              IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$RO_USER') THEN
                EXECUTE format('CREATE ROLE %I WITH LOGIN PASSWORD %L', '$RO_USER', gen_random_uuid()::text);
              END IF;
            END
            \$\$;
            GRANT CONNECT ON DATABASE "$PGDATABASE" TO "$RO_USER";
            GRANT USAGE ON SCHEMA public TO "$RO_USER";
            DO \$\$
            DECLARE
              t text;
            BEGIN
              FOREACH t IN ARRAY string_to_array('$RO_TABLES', ' ') LOOP
                IF to_regclass(format('public.%I', t)) IS NOT NULL THEN
                  EXECUTE format('GRANT SELECT ON TABLE public.%I TO %I', t, '$RO_USER');
                ELSE
                  RAISE WARNING 'table public.% does not exist yet; grant skipped (re-run after the owning service''s migrations)', t;
                END IF;
              END LOOP;
            END
            \$\$;
            SQL
            echo "Read-only role $RO_USER provisioned (requested SELECT on: $RO_TABLES)"
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
