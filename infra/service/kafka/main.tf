terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "~> 4.44"
    }
  }
}

ephemeral "vault_kv_secret_v2" "grafana_admin" {
  mount = "kv"
  name  = "grafana/admin/credentials"
}

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
}

data "vault_kv_secret_v2" "harbor_registry" {
  mount = "kv"
  name  = "harbor/users/harbor-devops/credentials"
}

data "kubernetes_service_v1" "harbor" {
  metadata {
    name      = "harbor"
    namespace = "harbor"
  }
}

resource "kubernetes_secret_v1" "harbor_push" {
  metadata {
    name      = "harbor-registry-secret"
    namespace = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  }
  type = "kubernetes.io/dockerconfigjson"
  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        (var.harbor_registry_host) = {
          username = data.vault_kv_secret_v2.harbor_registry.data["username"]
          password = data.vault_kv_secret_v2.harbor_registry.data["password"]
          auth     = base64encode("${data.vault_kv_secret_v2.harbor_registry.data["username"]}:${data.vault_kv_secret_v2.harbor_registry.data["password"]}")
        }
      }
    })
  }
}

resource "kubernetes_namespace_v1" "kafka_operator" {
  metadata {
    name = var.operator_namespace
  }
}

resource "kubernetes_namespace_v1" "kafka_broker" {
  metadata {
    name = var.broker_namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

resource "kubernetes_limit_range_v1" "resource_limits" {
  metadata {
    name      = "resource-limits"
    namespace = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "10m"
        memory = "32Mi"
      }
    }
  }
}

module "strimzi_operator" {
  source           = "./modules/strimzi-operator"
  namespace        = kubernetes_namespace_v1.kafka_operator.metadata[0].name
  watch_namespaces = [kubernetes_namespace_v1.kafka_broker.metadata[0].name]
  providers = {
    helm = helm
  }
}

module "strimzi_broker" {
  source             = "./modules/strimzi-broker"
  namespace          = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  kafka_cluster_name = var.kafka_cluster_name
  depends_on         = [module.strimzi_operator]
}

module "confluent_operator" {
  source           = "./modules/confluent-operator"
  namespace        = kubernetes_namespace_v1.kafka_operator.metadata[0].name
  watch_namespaces = [kubernetes_namespace_v1.kafka_broker.metadata[0].name]
  providers = {
    helm = helm
  }
}

module "confluent_schema_registry" {
  source             = "./modules/confluent-schema-registry"
  namespace          = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  kafka_cluster_name = var.kafka_cluster_name
  depends_on         = [module.confluent_operator]
}

module "kafka_connect" {
  source                   = "./modules/kafka-connect"
  namespace                = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  kafka_cluster_name       = var.kafka_cluster_name
  connect_cluster_name     = var.connect_cluster_name
  connect_image            = var.connect_image
  replicas                 = var.replicas
  image_pull_secret        = kubernetes_secret_v1.harbor_push.metadata[0].name
  additional_build_options = ["--insecure", "--insecure-pull", "--skip-tls-verify"]
  build_host_aliases = [
    {
      ip        = data.kubernetes_service_v1.harbor.spec[0].cluster_ip
      hostnames = [var.harbor_registry_host]
    },
  ]
  external_secret_mounts = [
    {
      name       = "catalog-postgres-debezium"
      mount_name = "debezium"
    },
  ]
  depends_on = [module.strimzi_broker, module.confluent_schema_registry]
}

module "debezium_connectors" {
  source                     = "./modules/debezium-connectors"
  namespace                  = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  kafka_cluster_name         = var.kafka_cluster_name
  connect_cluster_name       = var.connect_cluster_name
  catalog_postgres_namespace = var.catalog_postgres_namespace
  catalog_postgres_service   = var.catalog_postgres_service
  catalog_postgres_database  = var.catalog_postgres_database
  debezium_secret_namespace  = var.debezium_secret_namespace
  debezium_secret_name       = var.debezium_secret_name
  depends_on                 = [module.kafka_connect]
}

module "servicemonitor" {
  source             = "./modules/servicemonitor"
  namespace          = kubernetes_namespace_v1.kafka_broker.metadata[0].name
  kafka_cluster_name = var.kafka_cluster_name
  operator_namespace = var.operator_namespace
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "kafka"
  folder_title      = "Kafka"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "KafkaBrokerDown"
      expr        = "up{job=~\".*kafka.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Kafka broker is down"
      description = "Kafka broker {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "KafkaUnderReplicatedPartitions"
      expr        = "kafka_server_replicamanager_underreplicatedpartitions > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Kafka under-replicated partitions detected"
      description = "Kafka broker {{ $labels.instance }} has {{ $value }} under-replicated partitions."
    },
    {
      name        = "KafkaOfflinePartitions"
      expr        = "kafka_controller_kafkacontroller_offlinepartitionscount > 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Kafka offline partitions detected"
      description = "Kafka cluster has {{ $value }} offline partitions on {{ $labels.instance }}."
    },
    {
      name        = "KafkaConsumerLagHigh"
      expr        = "kafka_consumergroup_lag > 1000"
      for         = "15m"
      severity    = "warning"
      summary     = "Kafka consumer lag is high"
      description = "Kafka consumer group {{ $labels.consumergroup }} on topic {{ $labels.topic }} has a lag of {{ $value }}, exceeding 1000 threshold."
    },
    {
      name        = "KafkaBrokerDiskUsageHigh"
      expr        = "kafka_log_log_size / kafka_log_log_end_offset > 0.85"
      for         = "10m"
      severity    = "warning"
      summary     = "Kafka broker disk usage high"
      description = "Kafka broker {{ $labels.instance }} disk usage is {{ $value | humanizePercentage }}, exceeding 85% threshold."
    },
    {
      name        = "KafkaNoActiveController"
      expr        = "kafka_controller_kafkacontroller_activecontrollercount != 1"
      for         = "5m"
      severity    = "critical"
      summary     = "Kafka has no active controller"
      description = "Kafka cluster active controller count is {{ $value }}, expected exactly 1."
    },
    {
      name        = "KafkaISRShrinkRate"
      expr        = "rate(kafka_server_replicamanager_isrshrinks_total[5m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Kafka ISR shrink rate increasing"
      description = "Kafka broker {{ $labels.instance }} ISR shrink rate is {{ $value }}/s, indicating replica synchronization issues."
    },
    {
      name        = "KafkaRequestQueueOverflow"
      expr        = "kafka_network_requestchannel_request_queue_size > 100"
      for         = "5m"
      severity    = "warning"
      summary     = "Kafka request queue overflow"
      description = "Kafka broker {{ $labels.instance }} request queue size is {{ $value }}, exceeding 100 threshold."
    },
    {
      name        = "KafkaTopicPartitionSkew"
      expr        = "max without(partition) (kafka_log_log_size) / min without(partition) (kafka_log_log_size) > 1.5"
      for         = "10m"
      severity    = "warning"
      summary     = "Kafka topic partition skew detected"
      description = "Kafka topic {{ $labels.topic }} has a partition size skew ratio of {{ $value }}, exceeding 1.5 threshold."
    },
    {
      name        = "StrimziReconciliationFailed"
      expr        = "strimzi_reconciliations_failed_total > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Strimzi reconciliation failed"
      description = "Strimzi operator has {{ $value }} failed reconciliations for {{ $labels.kind }}/{{ $labels.name }}."
    },
    {
      name        = "StrimziReconciliationSlow"
      expr        = "strimzi_reconciliations_duration_seconds > 300"
      for         = "5m"
      severity    = "warning"
      summary     = "Strimzi reconciliation is slow"
      description = "Strimzi operator reconciliation for {{ $labels.kind }}/{{ $labels.name }} took {{ $value }}s, exceeding 300s threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
