locals {
  schema_name = "address-subject"
}

resource "kubernetes_manifest" "address_schema_config" {
  manifest = {
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata = {
      name      = local.schema_name
      namespace = var.namespace
    }
    data = {
      schema = jsonencode({
        namespace = "com.example"
        type      = "record"
        name      = "Address"
        fields = [
          {
            name = "street"
            type = "string"
          },
          {
            name    = "street2"
            type    = ["null", "string"]
            default = null
          },
          {
            name = "city"
            type = "string"
          },
          {
            name    = "pincode"
            type    = ["null", "int"]
            default = null
          }
        ]
      })
    }
  }
}

resource "kubernetes_manifest" "address_schema" {
  manifest = {
    apiVersion = "platform.confluent.io/v1beta1"
    kind       = "Schema"
    metadata = {
      name      = local.schema_name
      namespace = var.namespace
    }
    spec = {
      compatibilityLevel = "BACKWARD"
      name               = local.schema_name
      data = {
        configRef = kubernetes_manifest.address_schema_config.manifest.metadata.name
        format    = "avro"
      }
    }
  }
}
