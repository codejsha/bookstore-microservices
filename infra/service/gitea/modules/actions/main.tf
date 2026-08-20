terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

resource "null_resource" "runner_token_secret" {
  triggers = {
    secret_name = var.token_secret_name
    namespace   = var.namespace
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER  = var.admin_username
      GITEA_PASS  = var.admin_password
      API_URL     = var.gitea_api_url
      NAMESPACE   = var.namespace
      SECRET_NAME = var.token_secret_name
    }
    command = <<-EOT
      set -euo pipefail

      out=$(curl -s -w '\n%%{http_code}' -u "$GITEA_USER:$GITEA_PASS" \
        -X POST "$API_URL/admin/actions/runners/registration-token")
      code=$${out##*$'\n'}
      payload=$${out%$'\n'*}
      if [ "$code" != "200" ] && [ "$code" != "201" ]; then
        echo "gitea runner token: HTTP $code: $payload" >&2
        exit 1
      fi

      token=$(printf '%s' "$payload" | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')
      kubectl -n "$NAMESPACE" create secret generic "$SECRET_NAME" \
        --from-literal=token="$token" \
        --dry-run=client -o yaml | kubectl apply -f -
    EOT
  }
}

resource "helm_release" "actions" {
  namespace  = var.namespace
  name       = "gitea-actions"
  repository = "https://dl.gitea.com/charts"
  chart      = "actions"
  version    = "0.1.2"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 300

  set = [
    { name = "giteaRootURL", value = "http://${var.gitea_service_name}.${var.namespace}.svc.cluster.local:3000" },
    { name = "existingSecret", value = var.token_secret_name },
  ]

  depends_on = [null_resource.runner_token_secret]
}
