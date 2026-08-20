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
  }
}

data "vault_kv_secret_v2" "harbor_user" {
  mount = "kv"
  name  = "harbor/users/${var.harbor_user}/credentials"
}

data "vault_generic_secret" "cosign_key" {
  path = "transit/keys/${var.cosign_key_name}"
}

locals {
  transit_key       = jsondecode(data.vault_generic_secret.cosign_key.data_json)
  cosign_public_key = local.transit_key.keys[tostring(local.transit_key.latest_version)].public_key

  harbor_username = data.vault_kv_secret_v2.harbor_user.data["username"]
  harbor_password = data.vault_kv_secret_v2.harbor_user.data["password"]

  namespace_selector = {
    matchExpressions = [
      {
        key      = "kubernetes.io/metadata.name"
        operator = "In"
        values   = var.signed_namespaces
      }
    ]
  }

  resource_rules = [
    {
      apiGroups   = [""]
      apiVersions = ["v1"]
      operations  = ["CREATE", "UPDATE"]
      resources   = ["pods"]
    }
  ]

  attestor = {
    name = "releasekey"
    cosign = {
      key = {
        data = local.cosign_public_key
      }
      ctlog = {
        insecureIgnoreTlog = true
        insecureIgnoreSCT  = true
        url                = var.rekor_url
        rekorPubKey        = local.cosign_public_key
        ctLogPubKey        = local.cosign_public_key
      }
    }
  }

  credentials = {
    secrets = [var.image_pull_secret]
  }
}

resource "kubernetes_secret_v1" "harbor_pull" {
  metadata {
    name      = var.image_pull_secret
    namespace = var.namespace
  }
  type = "kubernetes.io/dockerconfigjson"
  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        (var.harbor_address) = {
          username = local.harbor_username
          password = local.harbor_password
          auth     = base64encode("${local.harbor_username}:${local.harbor_password}")
        }
      }
    })
  }
}

resource "kubernetes_manifest" "image_signature" {
  manifest = {
    apiVersion = "policies.kyverno.io/v1beta1"
    kind       = "ImageValidatingPolicy"
    metadata = {
      name = "bookstore-image-signature"
    }
    spec = {
      failurePolicy     = "Ignore"
      validationActions = [var.validation_action]
      evaluation = {
        admission  = { enabled = var.admission_enabled }
        background = { enabled = true }
      }
      matchConstraints = {
        namespaceSelector = local.namespace_selector
        resourceRules     = local.resource_rules
      }
      matchImageReferences = [
        { glob = var.image_glob }
      ]
      credentials = local.credentials
      attestors   = [local.attestor]
      validations = [
        {
          expression = "images.containers.map(image, verifyImageSignatures(image, [attestors.releasekey])).all(e, e > 0)"
          message    = "image is not signed by the bookstore release key"
        }
      ]
    }
  }

  field_manager {
    force_conflicts = true
  }
}

resource "kubernetes_manifest" "image_provenance" {
  count = var.provenance_policy_enabled ? 1 : 0

  manifest = {
    apiVersion = "policies.kyverno.io/v1beta1"
    kind       = "ImageValidatingPolicy"
    metadata = {
      name = "bookstore-image-provenance"
    }
    spec = {
      failurePolicy     = "Ignore"
      validationActions = [var.validation_action]
      evaluation = {
        admission  = { enabled = var.admission_enabled }
        background = { enabled = true }
      }
      matchConstraints = {
        namespaceSelector = local.namespace_selector
        resourceRules     = local.resource_rules
      }
      matchImageReferences = [
        { glob = var.image_glob }
      ]
      credentials = local.credentials
      attestors   = [local.attestor]
      attestations = [
        {
          name = "provenance"
          intoto = {
            type = var.provenance_predicate_type
          }
        }
      ]
      validations = [
        {
          expression = "images.containers.map(image, verifyAttestationSignatures(image, 'provenance', [attestors.releasekey])).all(e, e > 0)"
          message    = "image carries no build provenance signed by the bookstore release key"
        },
        {
          expression = "images.containers.map(image, extractPayload(image, 'provenance').predicate.runDetails.builder.id == '${var.builder_id}').all(e, e)"
          message    = "build provenance was not produced by the bookstore CI builder"
        }
      ]
    }
  }

  field_manager {
    force_conflicts = true
  }

  depends_on = [kubernetes_manifest.image_signature]
}
