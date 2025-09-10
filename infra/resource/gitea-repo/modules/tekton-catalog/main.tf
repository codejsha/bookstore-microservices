terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }
  }
}

resource "null_resource" "mirror" {
  triggers = {
    repo     = var.repo_name
    upstream = var.upstream_url
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER = var.gitea_username
      GITEA_PASS = var.gitea_password
      GITEA_AUTH = base64encode("${var.gitea_username}:${var.gitea_password}")
    }
    command = <<-EOT
      set -euo pipefail
      curl -s -u "$GITEA_USER:$GITEA_PASS" -X POST "${var.gitea_api_url}/orgs/${var.org_name}/repos" \
        -H 'Content-Type: application/json' \
        -d '{"name":"${var.repo_name}","private":true,"default_branch":"main"}' >/dev/null || true

      tmp="$(mktemp -d)"
      git clone --mirror "${var.upstream_url}" "$tmp/catalog.git"
      cd "$tmp/catalog.git"
      git for-each-ref --format='delete %(refname)' refs/pull | git update-ref --stdin
      git -c core.hooksPath=/dev/null -c http.extraheader="Authorization: Basic $GITEA_AUTH" \
        push --mirror "${var.gitea_http_url}/${var.org_name}/${var.repo_name}"
      cd - >/dev/null
      rm -rf "$tmp"
    EOT
  }
}
