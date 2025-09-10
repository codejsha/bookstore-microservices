resource "null_resource" "kustomize_build" {
  triggers = {
    manifests = sha1(join(",", [
      for f in sort(tolist(fileset(path.module, "**/*.yaml"))) :
      filesha1("${path.module}/${f}")
    ]))
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      set -euo pipefail
      kustomize build ${path.module} | kubectl apply --server-side --force-conflicts -f -
      kubectl wait --for=condition=Established crd/tektonconfigs.operator.tekton.dev --timeout=180s
      kubectl -n tekton-operator rollout status deploy/tekton-operator --timeout=300s
      kubectl -n tekton-operator rollout status deploy/tekton-operator-webhook --timeout=300s
    EOT
  }
}
