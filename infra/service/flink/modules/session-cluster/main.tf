resource "kubernetes_manifest" "flink_session_cluster" {
  manifest = {
    apiVersion = "flink.apache.org/v1beta1"
    kind       = "FlinkDeployment"
    metadata = {
      name      = "flink-session-cluster"
      namespace = var.namespace
    }
    spec = {
      image        = "flink:2.0"
      flinkVersion = "v2_0"
      flinkConfiguration = {
        "taskmanager.numberOfTaskSlots" = "2"
      }
      serviceAccount = "flink-operator"
      jobManager = {
        resource = {
          memory = "1024m"
          cpu    = 0.1
        }
      }
      taskManager = {
        replicas = 2
        resource = {
          memory = "1024m"
          cpu    = 0.1
        }
      }
    }
  }
}
