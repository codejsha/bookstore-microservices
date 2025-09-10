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

data "external" "monorepo_rev" {
  program = ["bash", "-c",
    "printf '{\"sha\":\"%s\"}' \"$(git -C '${var.content_dir}' rev-parse HEAD 2>/dev/null || echo none)\""
  ]
}

resource "null_resource" "push" {
  triggers = {
    content = sha1(join(",", [
      for f in sort(tolist(fileset(var.content_dir, "**"))) :
      filesha1("${var.content_dir}/${f}")
    ]))
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER   = var.gitea_username
      GITEA_PASS   = var.gitea_password
      GITEA_AUTH   = base64encode("${var.gitea_username}:${var.gitea_password}")
      REPO         = "${var.gitea_http_url}/${var.org_name}/${var.repo_name}.git"
      SRC          = var.content_dir
      MONOREPO_SHA = data.external.monorepo_rev.result.sha
    }
    command = <<-EOT
      set -euo pipefail
      curl -s -u "$GITEA_USER:$GITEA_PASS" -X POST "${var.gitea_api_url}/orgs/${var.org_name}/repos" \
        -H 'Content-Type: application/json' \
        -d '{"name":"${var.repo_name}","private":true,"default_branch":"main"}' >/dev/null || true

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
      mkdir -p task
      find "$SRC" -mindepth 1 -maxdepth 1 ! -name .git -exec cp -R {} task/ \;

      git_c add -A
      if git diff --cached --quiet; then
        echo "${var.repo_name}: catalog unchanged, nothing to push"
        exit 0
      fi
      git_c commit -q -m "catalog: sync ${var.repo_name} from monorepo" \
        -m "Monorepo-Commit: $MONOREPO_SHA"
      git_c push -q origin HEAD:main
    EOT
  }
}
