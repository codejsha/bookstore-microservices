terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }
    external = {
      source = "hashicorp/external"
    }
  }
}

data "external" "src_rev" {
  for_each = toset(var.services)
  program = ["bash", "-c",
    "printf '{\"sha\":\"%s\"}' \"$(git -C '${var.repo_root}' rev-parse 'HEAD:services/${each.key}' 2>/dev/null || echo none)\""
  ]
}

data "external" "monorepo_rev" {
  program = ["bash", "-c",
    "printf '{\"sha\":\"%s\"}' \"$(git -C '${var.repo_root}' rev-parse HEAD 2>/dev/null || echo none)\""
  ]
}

resource "null_resource" "source_push" {
  for_each = toset(var.services)

  triggers = {
    rev      = data.external.src_rev[each.key].result.sha
    repo     = "${each.key}${var.repo_suffix}"
    workflow = filemd5("${path.module}/files/tekton-pipeline.yml")
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER    = var.gitea_username
      GITEA_PASS    = var.gitea_password
      GITEA_AUTH    = base64encode("${var.gitea_username}:${var.gitea_password}")
      REPO          = "${var.gitea_http_url}/${var.org_name}/${each.key}${var.repo_suffix}.git"
      MONOREPO_SHA  = data.external.monorepo_rev.result.sha
      SUBTREE_SHA   = data.external.src_rev[each.key].result.sha
      SERVICE       = each.key
      WORKFLOW_FILE = abspath("${path.module}/files/tekton-pipeline.yml")
    }
    command = <<-EOT
      set -euo pipefail
      curl -s -u "$GITEA_USER:$GITEA_PASS" -X POST "${var.gitea_api_url}/orgs/${var.org_name}/repos" \
        -H 'Content-Type: application/json' \
        -d '{"name":"${each.key}${var.repo_suffix}","private":true,"default_branch":"main"}' >/dev/null || true

      git_c() { git -c core.hooksPath=/dev/null -c commit.gpgsign=false \
                    -c user.email=admin@example.com -c user.name=gitea_admin \
                    -c http.extraheader="Authorization: Basic $GITEA_AUTH" "$@"; }

      tmp="$(mktemp -d)"
      trap 'rm -rf "$tmp"' EXIT

      git_c clone -q "$REPO" "$tmp" 2>/dev/null || {
        git -c core.hooksPath=/dev/null init -q -b main "$tmp"
        git -C "$tmp" remote add origin "$REPO"
      }
      cd "$tmp"
      [ "$(git symbolic-ref -q --short HEAD || true)" = main ] || git checkout -q -B main

      find . -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
      git -C "${var.repo_root}" archive --format=tar "HEAD:services/$SERVICE" | tar -x -C .
      mkdir -p .gitea/workflows
      cp "$WORKFLOW_FILE" .gitea/workflows/tekton-pipeline.yml

      git_c add -A
      if git diff --cached --quiet; then
        echo "$SERVICE: source unchanged, nothing to push"
        exit 0
      fi
      git_c commit -q -m "sync $SERVICE from monorepo" \
        -m "Monorepo-Commit: $MONOREPO_SHA" \
        -m "Subtree-Sha: $SUBTREE_SHA"
      git_c push -q origin HEAD:main
    EOT
  }
}
