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

data "external" "head" {
  for_each = var.repos
  program = ["bash", "-c",
    "printf '{\"sha\":\"%s\"}' \"$(gh api 'repos/${var.github_owner}/${each.key}/commits/${var.github_branch}' --jq .sha 2>/dev/null || echo none)\""
  ]
}

resource "null_resource" "mirror" {
  for_each = var.repos

  triggers = {
    repo = each.key
    head = data.external.head[each.key].result.sha
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_AUTH = base64encode("${var.gitea_username}:${var.gitea_password}")
      GITEA_USER = var.gitea_username
      GITEA_PASS = var.gitea_password
    }
    command = <<-EOT
      set -euo pipefail

      gh_token="$(gh auth token)"
      if [ -z "$gh_token" ]; then
        echo "github-mirror: no GitHub token - run 'gh auth login' (the source repos are private)" >&2
        exit 1
      fi

      curl -s -u "$GITEA_USER:$GITEA_PASS" -X POST "${var.gitea_api_url}/orgs/${var.org_name}/repos" \
        -H 'Content-Type: application/json' \
        -d '{"name":"${each.key}","private":true,"default_branch":"${var.github_branch}"}' >/dev/null || true

      tmp="$(mktemp -d)"
      git clone --mirror \
        "https://x-access-token:$gh_token@github.com/${var.github_owner}/${each.key}.git" "$tmp/repo.git"
      cd "$tmp/repo.git"
      git for-each-ref --format='delete %(refname)' refs/pull | git update-ref --stdin
      git -c core.hooksPath=/dev/null -c http.extraheader="Authorization: Basic $GITEA_AUTH" \
        push --mirror "${var.gitea_http_url}/${var.org_name}/${each.key}"
      cd - >/dev/null
      rm -rf "$tmp"
    EOT
  }
}
