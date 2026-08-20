terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  catalog_cdc_tables = [
    "publisher",
    "author",
    "work",
    "edition",
    "subject",
    "work_author_mapping",
    "work_subject_mapping",
  ]
}

resource "kubernetes_manifest" "catalog_cdc_topic" {
  for_each = toset(local.catalog_cdc_tables)

  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaTopic"
    metadata = {
      name      = "catalog.public.${replace(each.key, "_", "-")}"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.kafka_cluster_name
      }
    }
    spec = {
      topicName  = "catalog.public.${each.key}"
      partitions = 3
      replicas   = 3
      config = {
        "retention.ms" = "-1"
      }
    }
  }
}

data "kubernetes_secret_v1" "catalog_debezium_src" {
  metadata {
    name      = var.debezium_secret_name
    namespace = var.debezium_secret_namespace
  }
}

resource "kubernetes_secret_v1" "catalog_debezium_mirror" {
  metadata {
    name      = var.debezium_secret_name
    namespace = var.namespace
  }
  type = "Opaque"
  data = data.kubernetes_secret_v1.catalog_debezium_src.data
}

resource "kubernetes_manifest" "catalog_debezium_connector" {
  depends_on = [
    kubernetes_secret_v1.catalog_debezium_mirror,
    kubernetes_manifest.catalog_cdc_topic,
  ]

  manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaConnector"
    metadata = {
      name      = "catalog-postgres-source"
      namespace = var.namespace
      labels = {
        "strimzi.io/cluster" = var.connect_cluster_name
      }
    }
    spec = {
      class    = "io.debezium.connector.postgresql.PostgresConnector"
      tasksMax = 1
      autoRestart = {
        enabled = true
      }
      config = {
        "database.hostname" = "${var.catalog_postgres_service}.${var.catalog_postgres_namespace}.svc.cluster.local"
        "database.port"     = 5432
        "database.user"     = "$${directory:/mnt/debezium:username}"
        "database.password" = "$${directory:/mnt/debezium:password}"
        "database.dbname"   = var.catalog_postgres_database

        "topic.prefix"                = "catalog"
        "plugin.name"                 = "pgoutput"
        "publication.name"            = "catalog_dbz_pub"
        "publication.autocreate.mode" = "filtered"
        "slot.name"                   = "catalog_dbz_slot"

        "schema.include.list" = "public"
        "table.include.list"  = join(",", [for t in local.catalog_cdc_tables : "public.${t}"])

        "snapshot.mode" = "initial"

        "decimal.handling.mode"        = "double"
        "time.precision.mode"          = "connect"
        "tombstones.on.delete"         = false
        "include.schema.changes"       = false
        "provide.transaction.metadata" = false
        "binary.handling.mode"         = "bytes"

        "key.converter"                       = "io.confluent.connect.avro.AvroConverter"
        "key.converter.schema.registry.url"   = "http://schemaregistry.${var.namespace}.svc.cluster.local:8081"
        "value.converter"                     = "io.confluent.connect.avro.AvroConverter"
        "value.converter.schema.registry.url" = "http://schemaregistry.${var.namespace}.svc.cluster.local:8081"
      }
    }
  }
}
