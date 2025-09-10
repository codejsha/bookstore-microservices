terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    gitea = {
      source  = "go-gitea/gitea"
      version = "~> 0.8"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
  }
}

provider "gitea" {
  base_url = "https://${var.gitea_address}"
  username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

data "vault_kv_secret_v2" "argocd_ci_token" {
  mount = "kv"
  name  = "argocd/dev-ci/token"
}

data "vault_kv_secret_v2" "gitea_webhook" {
  mount = "kv"
  name  = "gitea/webhook/credentials"
}

data "vault_kv_secret_v2" "gitea_resolver_token" {
  mount = "kv"
  name  = "gitea/resolver/token"
}

data "vault_kv_secret_v2" "gitea_admin" {
  mount = "kv"
  name  = "gitea/admin/credentials"
}

data "vault_kv_secret_v2" "nexus_publisher" {
  mount = "kv"
  name  = "nexus/publisher"
}

data "vault_kv_secret_v2" "harbor_ci" {
  mount = "kv"
  name  = "harbor/ci/credentials"
}

resource "kubernetes_secret_v1" "nexus_credentials" {
  metadata {
    name      = "nexus-credentials"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    username = data.vault_kv_secret_v2.nexus_publisher.data["username"]
    password = data.vault_kv_secret_v2.nexus_publisher.data["password"]
  }
}

locals {
  service_eventlisteners = merge(
    { for s in var.golang_services : s => "bookstore-golang" },
    { for s in var.gradle_services : s => "bookstore-gradle" },
    { for s in var.uv_services : s => "bookstore-uv" },
    { for s in var.vite_services : s => "bookstore-vite" },
  )
}

resource "kubernetes_resource_quota_v1" "ci_concurrency" {
  metadata {
    name      = "ci-concurrency"
    namespace = var.namespace
  }
  spec {
    hard = {
      pods = tostring(var.eventlistener_count + var.max_concurrent_pipelineruns * 2)
    }
  }
}

resource "kubernetes_secret_v1" "gitea_webhook" {
  metadata {
    name      = "gitea-webhook-secret"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    secretToken = data.vault_kv_secret_v2.gitea_webhook.data["token"]
  }
}

resource "vault_transit_secret_backend_key" "cosign" {
  backend                = "transit"
  name                   = "cosign-key"
  type                   = "ecdsa-p256"
  exportable             = false
  allow_plaintext_backup = false
  deletion_allowed       = true
}

resource "vault_policy" "tekton_chains_sign" {
  name   = "tekton-chains-sign"
  policy = <<-EOT
    path "transit/sign/${vault_transit_secret_backend_key.cosign.name}/*" {
      capabilities = ["update"]
    }
    path "transit/verify/${vault_transit_secret_backend_key.cosign.name}/*" {
      capabilities = ["update"]
    }
    path "transit/keys/${vault_transit_secret_backend_key.cosign.name}" {
      capabilities = ["read"]
    }
  EOT
}

resource "vault_kubernetes_auth_backend_role" "tekton_chains" {
  role_name                        = "tekton-chains-role"
  bound_service_account_names      = [var.chains_service_account]
  bound_service_account_namespaces = [var.chains_namespace]
  token_policies                   = [vault_policy.tekton_chains_sign.name]
  token_ttl                        = "3600"
}

resource "kubernetes_secret_v1" "gitea_resolver_token" {
  metadata {
    name      = var.resolver_token_secret
    namespace = var.resolver_namespace
  }
  type = "Opaque"
  data = {
    token = data.vault_kv_secret_v2.gitea_resolver_token.data["token"]
  }
}

resource "null_resource" "restart_resolver" {
  triggers = {
    token_secret = sha1(data.vault_kv_secret_v2.gitea_resolver_token.data["token"])
  }

  depends_on = [kubernetes_secret_v1.gitea_resolver_token]

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      set -euo pipefail
      kubectl -n ${var.resolver_namespace} delete pod -l app.kubernetes.io/name=resolvers
      kubectl -n ${var.resolver_namespace} rollout status deploy/tekton-pipelines-remote-resolvers --timeout=180s
    EOT
  }
}

resource "gitea_repository_actions_secret" "tekton_trigger_token" {
  for_each = local.service_eventlisteners

  repository_owner = var.gitea_org
  repository       = "${each.key}${var.source_repo_suffix}"
  secret_name      = "TEKTON_TRIGGER_TOKEN"
  secret_value     = data.vault_kv_secret_v2.gitea_webhook.data["token"]
}

resource "kubernetes_role_v1" "pipelinerun_reader" {
  metadata {
    name      = "pipelinerun-reader"
    namespace = var.namespace
  }
  rule {
    api_groups = ["tekton.dev"]
    resources  = ["pipelineruns"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_role_binding_v1" "pipelinerun_reader_default" {
  metadata {
    name      = "pipelinerun-reader-default"
    namespace = var.namespace
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = kubernetes_role_v1.pipelinerun_reader.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = "default"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "gitea_status_credentials" {
  metadata {
    name      = "gitea-status-credentials"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    username = data.vault_kv_secret_v2.gitea_admin.data["username"]
    password = data.vault_kv_secret_v2.gitea_admin.data["password"]
  }
}

resource "kubernetes_secret_v1" "harbor_ci_pull" {
  metadata {
    name      = var.registry_pull_secret
    namespace = var.namespace
  }
  type = "kubernetes.io/dockerconfigjson"
  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        for host in distinct([var.harbor_address, var.harbor_registry]) : host => {
          username = data.vault_kv_secret_v2.harbor_ci.data["username"]
          password = data.vault_kv_secret_v2.harbor_ci.data["password"]
          auth     = base64encode("${data.vault_kv_secret_v2.harbor_ci.data["username"]}:${data.vault_kv_secret_v2.harbor_ci.data["password"]}")
        }
      }
    })
  }
}

resource "kubernetes_default_service_account_v1" "ci_default" {
  metadata {
    namespace = var.namespace
  }
  image_pull_secret {
    name = kubernetes_secret_v1.harbor_ci_pull.metadata[0].name
  }
}

module "service_account" {
  for_each           = toset(var.service_accounts)
  source             = "./modules/serviceaccount"
  namespace          = var.namespace
  service_account    = each.key
  image_pull_secrets = [kubernetes_secret_v1.harbor_ci_pull.metadata[0].name]
}

module "trigger_binding" {
  source    = "./modules/triggerbinding"
  namespace = var.namespace
}

module "golang_template" {
  source                  = "./modules/golang-template"
  nexus_raw_url           = var.nexus_raw_url
  skip_codegen_verify     = var.skip_codegen_verify
  vault_internal_url      = var.vault_internal_url
  harbor_registry         = var.harbor_registry
  namespace               = var.namespace
  services                = var.golang_services
  argocd_internal_address = var.argocd_internal_address
  argocd_token            = data.vault_kv_secret_v2.argocd_ci_token.data["token"]
  harbor_address          = var.harbor_address
  harbor_username         = data.vault_kv_secret_v2.harbor_ci.data["username"]
  harbor_token            = data.vault_kv_secret_v2.harbor_ci.data["password"]
  vault_url               = var.vault_url
  nexus_docker_registry   = var.nexus_docker_registry
  nexus_username          = data.vault_kv_secret_v2.nexus_publisher.data["username"]
  nexus_password          = data.vault_kv_secret_v2.nexus_publisher.data["password"]
  providers = {
    vault = vault
  }
}

module "gradle_template" {
  source              = "./modules/gradle-template"
  nexus_maven_url     = var.nexus_maven_url
  nexus_raw_url       = var.nexus_raw_url
  skip_codegen_verify = var.skip_codegen_verify
  vault_internal_url  = var.vault_internal_url
  harbor_registry     = var.harbor_registry
  harbor_address      = var.harbor_address
  namespace           = var.namespace
  services            = var.gradle_services
  vault_url           = var.vault_url
  providers = {
    vault = vault
  }
}

module "uv_template" {
  source             = "./modules/uv-template"
  nexus_pypi_url     = var.nexus_pypi_url
  vault_internal_url = var.vault_internal_url
  harbor_registry    = var.harbor_registry
  namespace          = var.namespace
  services           = var.uv_services
  vault_url          = var.vault_url
  providers = {
    vault = vault
  }
}

module "vite_template" {
  source             = "./modules/vite-template"
  nexus_npm_url      = var.nexus_npm_url
  vault_internal_url = var.vault_internal_url
  harbor_registry    = var.harbor_registry
  namespace          = var.namespace
  services           = var.vite_services
  vault_url          = var.vault_url
  providers = {
    vault = vault
  }
}

module "gitea_webhook" {
  source             = "./modules/gitea-webhook"
  webhooks           = local.service_eventlisteners
  org_name           = var.gitea_org
  source_repo_suffix = var.source_repo_suffix
  gitea_api_url      = "https://${var.gitea_address}/api/v1"
  gitea_username     = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password     = data.vault_kv_secret_v2.gitea_admin.data["password"]

  depends_on = [
    module.golang_template,
    module.gradle_template,
    module.uv_template,
    module.vite_template,
  ]
}
