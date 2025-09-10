terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  sql_dir = var.sql_dir
}

resource "kubernetes_config_map_v1" "catalog_indexer_sql" {
  metadata {
    name      = "catalog-indexer-sql"
    namespace = var.namespace
  }
  data = {
    "01-sources.sql"  = file("${local.sql_dir}/01-sources.sql")
    "02-sink.sql"     = file("${local.sql_dir}/02-sink.sql")
    "03-pipeline.sql" = file("${local.sql_dir}/03-pipeline.sql")
  }
}

resource "kubernetes_secret_v1" "catalog_indexer_opensearch" {
  metadata {
    name      = var.opensearch_secret_name
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    username = var.opensearch_username
    password = var.opensearch_password
  }
}

resource "kubernetes_secret_v1" "catalog_indexer_s3" {
  metadata {
    name      = "catalog-indexer-s3"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    access-key-id     = var.s3_access_key_id
    secret-access-key = var.s3_secret_key
  }
}

resource "kubernetes_secret_v1" "harbor_pull" {
  metadata {
    name      = var.image_pull_secret_name
    namespace = var.namespace
  }
  type = "kubernetes.io/dockerconfigjson"
  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        (var.registry_host) = {
          username = var.harbor_pull_username
          password = var.harbor_pull_password
          auth     = base64encode(format("%s:%s", var.harbor_pull_username, var.harbor_pull_password))
        }
      }
    })
  }
}

resource "kubernetes_manifest" "catalog_indexer" {
  depends_on = [
    kubernetes_config_map_v1.catalog_indexer_sql,
    kubernetes_secret_v1.catalog_indexer_opensearch,
    kubernetes_secret_v1.catalog_indexer_s3,
    kubernetes_secret_v1.harbor_pull,
  ]

  manifest = {
    apiVersion = "flink.apache.org/v1beta1"
    kind       = "FlinkDeployment"
    metadata = {
      name      = "catalog-indexer"
      namespace = var.namespace
    }
    spec = {
      image          = var.indexer_image
      flinkVersion   = "v1_20"
      serviceAccount = var.flink_service_account

      flinkConfiguration = {
        "taskmanager.numberOfTaskSlots"                        = "2"
        "execution.checkpointing.interval"                     = "120s"
        "execution.checkpointing.mode"                         = "EXACTLY_ONCE"
        "execution.checkpointing.timeout"                      = "30min"
        "execution.checkpointing.tolerable-failed-checkpoints" = "10"
        "execution.checkpointing.unaligned.enabled"            = "true"

        "state.backend.type"        = "rocksdb"
        "state.backend.incremental" = "true"

        "state.checkpoints.dir" = "s3://${var.s3_bucket}/catalog-indexer/checkpoints"
        "state.savepoints.dir"  = "s3://${var.s3_bucket}/catalog-indexer/savepoints"
        "s3.endpoint"           = var.s3_endpoint
        "s3.path.style.access"  = "true"

        "table.optimizer.non-deterministic-update.strategy" = "TRY_RESOLVE"
      }

      podTemplate = {
        apiVersion = "v1"
        kind       = "Pod"
        metadata = {
          name = "catalog-indexer"
          annotations = {
            "bookstore.dev/sql-checksum" = sha256(join("", values(kubernetes_config_map_v1.catalog_indexer_sql.data)))
          }
        }
        spec = {
          imagePullSecrets = [
            { name = kubernetes_secret_v1.harbor_pull.metadata[0].name }
          ]
          containers = [
            {
              name = "flink-main-container"
              env = [
                {
                  name = "OPENSEARCH_USERNAME"
                  valueFrom = {
                    secretKeyRef = {
                      name = kubernetes_secret_v1.catalog_indexer_opensearch.metadata[0].name
                      key  = "username"
                    }
                  }
                },
                {
                  name = "OPENSEARCH_PASSWORD"
                  valueFrom = {
                    secretKeyRef = {
                      name = kubernetes_secret_v1.catalog_indexer_opensearch.metadata[0].name
                      key  = "password"
                    }
                  }
                },
                {
                  name = "AWS_ACCESS_KEY_ID"
                  valueFrom = {
                    secretKeyRef = {
                      name = kubernetes_secret_v1.catalog_indexer_s3.metadata[0].name
                      key  = "access-key-id"
                    }
                  }
                },
                {
                  name = "AWS_SECRET_ACCESS_KEY"
                  valueFrom = {
                    secretKeyRef = {
                      name = kubernetes_secret_v1.catalog_indexer_s3.metadata[0].name
                      key  = "secret-access-key"
                    }
                  }
                }
              ]
              volumeMounts = [
                {
                  name      = "catalog-indexer-sql"
                  mountPath = "/opt/flink/sql"
                }
              ]
            }
          ]
          volumes = [
            {
              name = "catalog-indexer-sql"
              configMap = {
                name = kubernetes_config_map_v1.catalog_indexer_sql.metadata[0].name
              }
            }
          ]
        }
      }

      jobManager = {
        resource = {
          memory = "1024m"
          cpu    = 0.5
        }
      }
      taskManager = {
        resource = {
          memory = "2048m"
          cpu    = 0.5
        }
      }

      job = {
        jarURI      = "local:///opt/flink/lib/flink-sql-runner.jar"
        parallelism = 2
        upgradeMode = "stateless"
        state       = "running"
        args = [
          "/opt/flink/sql/01-sources.sql",
          "/opt/flink/sql/02-sink.sql",
          "/opt/flink/sql/03-pipeline.sql",
        ]
      }
    }
  }
}
