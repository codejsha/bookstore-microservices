resource "kubernetes_manifest" "git_triggerbinding" {
  manifest = {
    apiVersion = "triggers.tekton.dev/v1alpha1"
    kind       = "TriggerBinding"
    metadata = {
      name      = "git-triggerbinding"
      namespace = var.namespace
    }
    spec = {
      params = [
        {
          name  = "service-name"
          value = "$(extensions.service.name)"
        },
        {
          name  = "git-repo-url"
          value = "$(extensions.repo.url)"
        },
        {
          name  = "git-repo-name"
          value = "$(extensions.repo.name)"
        },
        {
          name  = "git-ref"
          value = "$(extensions.ref)"
        },
        {
          name  = "git-revision"
          value = "$(extensions.after)"
        },
      ]
    }
  }
}
