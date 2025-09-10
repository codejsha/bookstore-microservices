resource "kubernetes_manifest" "kafka_podmonitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "kafka-broker"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
        release              = "prometheus"
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "strimzi.io/cluster" = var.kafka_cluster_name
          "strimzi.io/kind"    = "Kafka"
        }
      }
      podTargetLabels = [
        "strimzi.io/cluster",
        "strimzi.io/name",
        "strimzi.io/kind",
      ]
      podMetricsEndpoints = [
        {
          path = "/metrics"
          port = "tcp-prometheus"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "kafka_exporter_podmonitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "kafka-exporter"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
        release              = "prometheus"
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "strimzi.io/cluster" = var.kafka_cluster_name
          "strimzi.io/name"    = "${var.kafka_cluster_name}-kafka-exporter"
        }
      }
      podTargetLabels = [
        "strimzi.io/cluster",
        "strimzi.io/name",
        "strimzi.io/kind",
      ]
      podMetricsEndpoints = [
        {
          path = "/metrics"
          port = "tcp-prometheus"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "kafka_connect_podmonitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "kafka-connect"
      namespace = var.namespace
      labels = {
        "strimzi.io/kind" = "KafkaConnect"
        release           = "prometheus"
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "strimzi.io/kind" = "KafkaConnect"
        }
      }
      podTargetLabels = [
        "strimzi.io/cluster",
        "strimzi.io/name",
        "strimzi.io/kind",
      ]
      podMetricsEndpoints = [
        {
          path = "/metrics"
          port = "tcp-prometheus"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "strimzi_operator_podmonitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "strimzi-cluster-operator"
      namespace = var.operator_namespace
      labels = {
        "strimzi.io/kind" = "cluster-operator"
        release           = "prometheus"
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "strimzi.io/kind" = "cluster-operator"
        }
      }
      podMetricsEndpoints = [
        {
          path = "/metrics"
          port = "http"
        }
      ]
    }
  }
}
