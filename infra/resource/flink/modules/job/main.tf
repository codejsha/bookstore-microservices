resource "kubernetes_manifest" "flink_session_job" {
  manifest = {
    apiVersion = "flink.apache.org/v1beta1"
    kind       = "FlinkSessionJob"
    metadata = {
      name      = var.job_name
      namespace = var.namespace
    }
    spec = {
      deploymentName = var.session_cluster_name
      job = {
        jarURI      = var.job_jar_uri
        parallelism = var.job_parallelism
        upgradeMode = "stateless"
        args        = var.job_args
      }
    }
  }
}
