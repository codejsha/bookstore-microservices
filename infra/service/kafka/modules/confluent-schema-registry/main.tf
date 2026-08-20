resource "kubernetes_manifest" "kafka_topic_schemas" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaTopic"
    metadata = {
      name      = "schemas-metadata"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
      }
    }
    spec = {
      topicName  = "schemas.metadata"
      partitions = 3
      replicas   = 3
      config = {
        "cleanup.policy" = "compact"
      }
    }
  }
}

resource "kubernetes_manifest" "schema_registry" {
  depends_on = [kubernetes_manifest.kafka_topic_schemas]
  manifest = {
    apiVersion = "platform.confluent.io/v1beta1"
    kind       = "SchemaRegistry"
    metadata = {
      name      = "schemaregistry"
      namespace = var.namespace
    }
    spec = {
      replicas = 1
      image = {
        application = "confluentinc/cp-schema-registry:8.1.0"
        init        = "confluentinc/confluent-init-container:3.1.0"
      }
      dependencies = {
        kafka = {
          bootstrapEndpoint = "${var.kafka_cluster_name}-kafka-bootstrap.${var.namespace}.svc.cluster.local:9093"
        }
      }
    }
  }
}
