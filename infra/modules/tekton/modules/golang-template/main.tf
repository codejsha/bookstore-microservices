resource "kubernetes_persistent_volume_claim" "tekton_shared_workspace" {
  metadata {
    name = "tekton-shared-workspace"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "5Gi"
      }
    }
    storage_class_name = "standard"
  }
}

locals {
  admin_service     = "admin"
  catalog_service   = "catalog"
  customer_service  = "customer"
  identity_service  = "identity"
  inventory_service = "inventory"
  order_service     = "order"
  payment_service   = "payment"
}

resource "kubernetes_manifest" "repo_ssh_auth" {
  for_each = {
    for k, v in {
      (local.admin_service)     = data.vault_generic_secret.admin_helm_ssh_keys.data
      (local.catalog_service)   = data.vault_generic_secret.catalog_helm_ssh_keys.data
      (local.customer_service)  = data.vault_generic_secret.customer_helm_ssh_keys.data
      (local.identity_service)  = data.vault_generic_secret.identity_helm_ssh_keys.data
      (local.inventory_service) = data.vault_generic_secret.inventory_helm_ssh_keys.data
      (local.order_service)     = data.vault_generic_secret.order_helm_ssh_keys.data
      (local.payment_service)   = data.vault_generic_secret.payment_helm_ssh_keys.data
    } : k => v

  }
  manifest = yamldecode(templatefile("${path.module}/manifests/repo-ssh-auth.yaml", {
    name       = "${each.key}-repo-ssh-auth"
    namespace  = var.namespace
    gitea_fqdn = var.gitea_fqdn
    ssh_private_key = base64encode(chomp(each.value["private"])) # Encode it
  }))
}

resource "kubernetes_manifest" "docker_configmap" {
  manifest = yamldecode(templatefile("${path.module}/manifests/docker-configmap.yaml", {
    namespace   = var.namespace,
    harbor_fqdn = var.harbor_fqdn
    harbor_auth = base64encode("${var.harbor_username}:${var.harbor_token}")
  }))
}

resource "kubernetes_manifest" "argocd_configmap" {
  manifest = yamldecode(templatefile("${path.module}/manifests/argocd-configmap.yaml", {
    namespace = var.namespace,
    argocd_address = base64encode(var.argocd_address)
  }))
}

resource "kubernetes_manifest" "argocd_secret" {
  manifest = yamldecode(templatefile("${path.module}/manifests/argocd-secret.yaml", {
    namespace = var.namespace,
    argocd_token = base64encode(var.argocd_token)
  }))
}

resource "kubernetes_manifest" "tekton_trigger_template" {
  manifest = yamldecode(templatefile("${path.module}/manifests/triggertemplate.yaml", {
    namespace = var.namespace,
  }))
}

resource "null_resource" "tekton_event_listener" {
  depends_on = [kubernetes_manifest.tekton_trigger_template]
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "kubectl apply -n ${var.namespace} -f ${path.module}/manifests/eventlistener.yaml"
  }
}
