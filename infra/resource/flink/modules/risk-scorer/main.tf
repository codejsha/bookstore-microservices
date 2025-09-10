terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_manifest" "access_log_topic" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaTopic"
    metadata = {
      name      = "bookstore.access-log"
      namespace = var.kafka_namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
      }
    }
    spec = {
      topicName  = "bookstore.access-log"
      partitions = 3
      replicas   = 3
      config = {
        "retention.ms" = "3600000"
      }
    }
  }
}

resource "kubernetes_config_map_v1" "risk_scorer_sql" {
  metadata {
    name      = "risk-scorer-sql"
    namespace = var.namespace
  }
  data = {
    "01-source.sql"   = file("${var.sql_dir}/01-source.sql")
    "02-sink.sql"     = file("${var.sql_dir}/02-sink.sql")
    "03-pipeline.sql" = file("${var.sql_dir}/03-pipeline.sql")
  }
}

resource "kubernetes_manifest" "risk_scorer" {
  depends_on = [
    kubernetes_config_map_v1.risk_scorer_sql,
  ]

  manifest = {
    apiVersion = "flink.apache.org/v1beta1"
    kind       = "FlinkDeployment"
    metadata = {
      name      = "risk-scorer"
      namespace = var.namespace
    }
    spec = {
      image          = var.runner_image
      flinkVersion   = "v1_20"
      serviceAccount = var.flink_service_account

      flinkConfiguration = {
        "taskmanager.numberOfTaskSlots"                        = "2"
        "execution.checkpointing.interval"                     = "60s"
        "execution.checkpointing.mode"                         = "AT_LEAST_ONCE"
        "execution.checkpointing.timeout"                      = "10min"
        "execution.checkpointing.tolerable-failed-checkpoints" = "10"

        "state.backend.type"        = "rocksdb"
        "state.backend.incremental" = "true"

        "state.checkpoints.dir" = "s3://${var.s3_bucket}/risk-scorer/checkpoints"
        "state.savepoints.dir"  = "s3://${var.s3_bucket}/risk-scorer/savepoints"
        "s3.endpoint"           = var.s3_endpoint
        "s3.path.style.access"  = "true"

        "table.local-time-zone" = "UTC"
      }

      podTemplate = {
        apiVersion = "v1"
        kind       = "Pod"
        metadata = {
          name = "risk-scorer"
          annotations = {
            "bookstore.dev/sql-checksum" = sha256(join("", values(kubernetes_config_map_v1.risk_scorer_sql.data)))
          }
        }
        spec = {
          imagePullSecrets = [
            { name = var.image_pull_secret_name }
          ]
          containers = [
            {
              name = "flink-main-container"
              env = [
                {
                  name  = "VALKEY_HOST"
                  value = var.valkey_host
                },
                {
                  name  = "RISK_SCORE_THRESHOLD"
                  value = tostring(var.risk_score_threshold)
                },
                {
                  name  = "AUTO_FLAG_TTL_SECONDS"
                  value = tostring(var.auto_flag_ttl_seconds)
                },
                {
                  name = "AWS_ACCESS_KEY_ID"
                  valueFrom = {
                    secretKeyRef = {
                      name = var.s3_secret_name
                      key  = "access-key-id"
                    }
                  }
                },
                {
                  name = "AWS_SECRET_ACCESS_KEY"
                  valueFrom = {
                    secretKeyRef = {
                      name = var.s3_secret_name
                      key  = "secret-access-key"
                    }
                  }
                }
              ]
              volumeMounts = [
                {
                  name      = "risk-scorer-sql"
                  mountPath = "/opt/flink/sql"
                }
              ]
            }
          ]
          volumes = [
            {
              name = "risk-scorer-sql"
              configMap = {
                name = kubernetes_config_map_v1.risk_scorer_sql.metadata[0].name
              }
            }
          ]
        }
      }

      jobManager = {
        resource = {
          memory = "1024m"
          cpu    = 0.25
        }
      }
      taskManager = {
        resource = {
          memory = "1536m"
          cpu    = 0.25
        }
      }

      job = {
        jarURI      = "local:///opt/flink/lib/flink-sql-runner.jar"
        parallelism = 1
        upgradeMode = "stateless"
        state       = "running"
        args = [
          "/opt/flink/sql/01-source.sql",
          "/opt/flink/sql/02-sink.sql",
          "/opt/flink/sql/03-pipeline.sql",
        ]
      }
    }
  }
}
