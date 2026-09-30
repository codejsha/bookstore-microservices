resource "null_resource" "kustomize_build" {
  triggers = {
    kustomization = filesha256("${path.module}/kustomization.yaml")
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      set -euo pipefail
      kustomize build ${path.module} | kubectl apply --server-side --force-conflicts -f -
      kubectl wait --for=condition=Established --timeout=180s \
        crd/flinkdeployments.flink.apache.org \
        crd/flinkbluegreendeployments.flink.apache.org \
        crd/flinksessionjobs.flink.apache.org \
        crd/flinkstatesnapshots.flink.apache.org
    EOT
  }
}
