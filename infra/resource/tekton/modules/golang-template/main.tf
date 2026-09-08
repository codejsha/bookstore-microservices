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

data "vault_kv_secret_v2" "lib_ssh_key" {
  mount = "kv"
  name  = "gitea/ssh/${var.lib_repo}"
}

resource "kubernetes_secret_v1" "lib_repo_ssh_auth" {
  metadata {
    name      = "${var.lib_repo}-repo-ssh-auth"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    "id_rsa" = chomp(data.vault_kv_secret_v2.lib_ssh_key.data["private"])
  }
}

resource "kubernetes_secret_v1" "dockerconfig" {
  metadata {
    name      = "dockerconfig-secret"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    "config.json" = jsonencode({
      auths = merge(
        {
          for host in distinct([var.harbor_address, var.harbor_registry]) : host => {
            auth = base64encode("${var.harbor_username}:${var.harbor_token}")
          }
        },
        {
          (var.nexus_docker_registry) = {
            auth = base64encode("${var.nexus_username}:${var.nexus_password}")
          }
        },
      )
    })
  }
}

resource "kubernetes_config_map_v1" "argocd_env" {
  metadata {
    name      = "argocd-env-configmap"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"    = "argocd-env-configmap"
      "app.kubernetes.io/part-of" = "argocd"
    }
  }
  data = {
    ARGOCD_SERVER = "${var.argocd_internal_address}:80"
    ARGOCD_OPTS   = "--plaintext"
  }
}

resource "kubernetes_secret_v1" "argocd_env" {
  metadata {
    name      = "argocd-env-secret"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"    = "argocd-env-secret"
      "app.kubernetes.io/part-of" = "argocd"
    }
  }
  type = "Opaque"
  data = {
    ARGOCD_AUTH_TOKEN = var.argocd_token
  }
}

resource "kubernetes_manifest" "tekton_trigger_template" {
  manifest = yamldecode(templatefile("${path.module}/manifests/triggertemplate.yaml", {
    namespace              = var.namespace,
    vault_url              = var.vault_url,
    harbor_registry        = var.harbor_registry,
    harbor_address         = var.harbor_address,
    vault_internal_url     = var.vault_internal_url,
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
