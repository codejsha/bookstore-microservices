terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

data "vault_kv_secret_v2" "source_ssh_keys" {
  for_each = toset(var.services)
  mount    = "kv"
  name     = "gitea/ssh/${each.key}-source"
}

data "vault_kv_secret_v2" "helm_ssh_keys" {
  for_each = toset(var.services)
  mount    = "kv"
  name     = "gitea/ssh/${each.key}-helm"
}

resource "kubernetes_secret_v1" "source_repo_ssh_auth" {
  for_each = toset(var.services)
  metadata {
    name      = "${each.key}-repo-ssh-auth"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    "id_rsa" = chomp(data.vault_kv_secret_v2.source_ssh_keys[each.key].data["private"])
  }
}

resource "kubernetes_secret_v1" "helm_repo_ssh_auth" {
  for_each = toset(var.services)
  metadata {
    name      = "${each.key}-helm-repo-ssh-auth"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    "id_rsa" = chomp(data.vault_kv_secret_v2.helm_ssh_keys[each.key].data["private"])
  }
}

resource "kubernetes_manifest" "tekton_trigger_template" {
  manifest = yamldecode(templatefile("${path.module}/manifests/triggertemplate.yaml", {
    namespace              = var.namespace,
    vault_url              = var.vault_url,
    harbor_registry        = var.harbor_registry,
    vault_internal_url     = var.vault_internal_url,
    nexus_pypi_url         = var.nexus_pypi_url,
    nexus_raw_url          = var.nexus_raw_url,
    codegen_verify_targets = var.codegen_verify_targets,
  }))

  field_manager {
    force_conflicts = true
  }
}

resource "kubernetes_manifest" "tekton_cd_trigger_template" {
  manifest = yamldecode(templatefile("${path.module}/manifests/cd-triggertemplate.yaml", {
    namespace          = var.namespace,
    harbor_registry    = var.harbor_registry,
    vault_internal_url = var.vault_internal_url,
  }))

  field_manager {
    force_conflicts = true
  }
}

resource "null_resource" "tekton_event_listener" {
  depends_on = [
    kubernetes_manifest.tekton_trigger_template,
    kubernetes_manifest.tekton_cd_trigger_template,
  ]
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "kubectl apply -n ${var.namespace} -f ${path.module}/manifests/eventlistener.yaml"
  }
}
