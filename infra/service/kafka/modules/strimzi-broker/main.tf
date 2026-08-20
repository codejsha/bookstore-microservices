resource "kubernetes_manifest" "kafka_node_pool_controller" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaNodePool"
    metadata = {
      name      = "controller"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
      }
    }
    spec = {
      replicas = 3
      roles    = ["controller"]
      resources = {
        requests = {
          cpu    = "10m"
          memory = "384Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "768Mi"
        }
      }
      storage = {
        type = "jbod"
        volumes = [
          {
            id            = 0
            type          = "persistent-claim"
            size          = "100Gi"
            kraftMetadata = "shared"
          }
        ]
      }
    }
  }
}

resource "kubernetes_manifest" "kafka_node_pool_broker" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaNodePool"
    metadata = {
      name      = "broker"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
      }
    }
    spec = {
      replicas = 3
      roles    = ["broker"]
      resources = {
        requests = {
          cpu    = "10m"
          memory = "384Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "768Mi"
        }
      }
      storage = {
        type = "jbod"
        volumes = [
          {
            id            = 0
            type          = "persistent-claim"
            size          = "100Gi"
            kraftMetadata = "shared"
          }
        ]
      }
    }
  }
}

resource "kubernetes_config_map_v1" "kafka_jmx_metrics" {
  metadata {
    name      = "kafka-metrics"
    namespace = var.namespace
    labels = {
      app = "strimzi"
    }
  }
  data = {
    "kafka-metrics-config.yml" = file("${path.module}/kafka-metrics-config.yml")
  }
}

resource "kubernetes_manifest" "kafka_cluster" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "Kafka"
    metadata = {
      name      = var.kafka_cluster_name
      namespace = var.namespace
    }
    spec = {
      kafka = {
        version         = "4.2.0"
        metadataVersion = "4.2-IV1"
        listeners = [
          {
            name = "internal"
            port = 9093
            type = "internal"
            tls  = false
          }
        ]
        config = {
          "default.replication.factor"               = 3
          "min.insync.replicas"                      = 2
          "offsets.topic.replication.factor"         = 3
          "transaction.state.log.replication.factor" = 3
          "transaction.state.log.min.isr"            = 2
          "auto.create.topics.enable"                = "false"
        }
        metricsConfig = {
          type = "jmxPrometheusExporter"
          valueFrom = {
            configMapKeyRef = {
              name = kubernetes_config_map_v1.kafka_jmx_metrics.metadata[0].name
              key  = "kafka-metrics-config.yml"
            }
          }
        }
      }
      entityOperator = {
        topicOperator = {
          resources = {
            requests = { cpu = "10m", memory = "256Mi" }
            limits   = { cpu = "500m", memory = "512Mi" }
          }
        }
        userOperator = {
          resources = {
            requests = { cpu = "10m", memory = "256Mi" }
            limits   = { cpu = "500m", memory = "512Mi" }
          }
        }
      }
      kafkaExporter = {
        topicRegex = ".*"
        groupRegex = ".*"
      }
    }
  }
  depends_on = [kubernetes_config_map_v1.kafka_jmx_metrics]
}
