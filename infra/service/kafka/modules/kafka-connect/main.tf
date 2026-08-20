locals {
  external_volumes = [
    for s in var.external_secret_mounts : {
      name   = s.mount_name
      secret = { secretName = s.name }
    }
  ]
  external_volume_mounts = [
    for s in var.external_secret_mounts : {
      name      = s.mount_name
      mountPath = "/mnt/${s.mount_name}"
    }
  ]
  pod_template = merge(
    var.image_pull_secret != "" ? {
      imagePullSecrets = [{ name = var.image_pull_secret }]
    } : {},
    length(local.external_volumes) > 0 ? {
      volumes = local.external_volumes
    } : {},
  )
  connect_template = merge(
    length(var.build_host_aliases) > 0 ? {
      buildPod = {
        hostAliases = var.build_host_aliases
      }
    } : {},
    length(local.pod_template) > 0 ? {
      pod = local.pod_template
    } : {},
    length(local.external_volume_mounts) > 0 ? {
      connectContainer = {
        volumeMounts = local.external_volume_mounts
      }
    } : {},
  )
  template_block = length(local.connect_template) > 0 ? { template = local.connect_template } : {}
}

resource "kubernetes_manifest" "kafka_connect" {
  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaConnect"
    metadata = {
      name      = var.connect_cluster_name
      namespace = var.namespace
      annotations = {
        "strimzi.io/use-connector-resources" = "true"
      }
    }
    spec = merge({
      version          = "4.2.0"
      replicas         = var.replicas
      bootstrapServers = "${var.kafka_cluster_name}-kafka-bootstrap.${var.namespace}.svc.cluster.local:9093"

      groupId            = "${var.connect_cluster_name}-cluster-group"
      offsetStorageTopic = "${var.connect_cluster_name}-cluster-offsets"
      configStorageTopic = "${var.connect_cluster_name}-cluster-configs"
      statusStorageTopic = "${var.connect_cluster_name}-cluster-status"

      config = {
        "config.providers"                    = "directory"
        "config.providers.directory.class"    = "org.apache.kafka.common.config.provider.DirectoryConfigProvider"
        "offset.storage.replication.factor"   = 3
        "config.storage.replication.factor"   = 3
        "status.storage.replication.factor"   = 3
        "key.converter"                       = "io.confluent.connect.avro.AvroConverter"
        "value.converter"                     = "io.confluent.connect.avro.AvroConverter"
        "key.converter.schema.registry.url"   = "http://schemaregistry.${var.namespace}.svc.cluster.local:8081"
        "value.converter.schema.registry.url" = "http://schemaregistry.${var.namespace}.svc.cluster.local:8081"
        "key.converter.schemas.enable"        = true
        "value.converter.schemas.enable"      = true
      }

      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }

      build = {
        output = merge(
          {
            type  = "docker"
            image = var.connect_image
          },
          var.image_pull_secret != "" ? {
            pushSecret = var.image_pull_secret
          } : {},
          length(var.additional_build_options) > 0 ? {
            additionalBuildOptions = var.additional_build_options
          } : {}
        )
        plugins = [
          {
            name = "debezium-postgres-connector"
            artifacts = [
              {
                type = "tgz"
                url  = "https://repo1.maven.org/maven2/io/debezium/debezium-connector-postgres/3.0.7.Final/debezium-connector-postgres-3.0.7.Final-plugin.tar.gz"
              }
            ]
          },
          {
            name = "confluent-avro-converter"
            artifacts = [
              {
                type = "zip"
                url  = "https://hub-downloads.confluent.io/api/plugins/confluentinc/kafka-connect-avro-converter/versions/7.7.1/confluentinc-kafka-connect-avro-converter-7.7.1.zip"
              }
            ]
          }
        ]
      }

      logging = {
        type = "inline"
        loggers = {
          "rootLogger.level" = "INFO"
        }
      }

      metricsConfig = {
        type = "jmxPrometheusExporter"
        valueFrom = {
          configMapKeyRef = {
            name = "${var.connect_cluster_name}-metrics"
            key  = "metrics-config.yml"
          }
        }
      }
    }, local.template_block)
  }
}

resource "kubernetes_config_map_v1" "connect_metrics" {
  metadata {
    name      = "${var.connect_cluster_name}-metrics"
    namespace = var.namespace
  }
  data = {
    "metrics-config.yml" = file("${path.module}/connect-metrics-config.yml")
  }
}
