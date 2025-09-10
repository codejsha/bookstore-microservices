terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }
  }
}

locals {
  charts = { for r in var.app_repos : r => replace(r, "-helm", "") }
}

resource "null_resource" "chart_push" {
  for_each = local.charts

  triggers = {
    content = sha1(join(",", [
      for f in sort(tolist(fileset("${var.helm_dir}/${each.value}", "**"))) :
      filesha1("${var.helm_dir}/${each.value}/${f}")
    ]))
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER = var.gitea_username
      GITEA_PASS = var.gitea_password
      GITEA_AUTH = base64encode("${var.gitea_username}:${var.gitea_password}")
      SRC        = "${var.helm_dir}/${each.value}"
      REPO       = "${var.gitea_http_url}/${var.org_name}/${each.key}.git"
    }
    command = <<-EOT
      set -euo pipefail
      curl -s -u "$GITEA_USER:$GITEA_PASS" -X POST "${var.gitea_api_url}/orgs/${var.org_name}/repos" \
        -H 'Content-Type: application/json' \
        -d '{"name":"${each.key}","private":true,"default_branch":"main"}' >/dev/null || true

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

      tag=""
      if [ -f values.yaml ]; then
        tag="$(yq -r '.image.tag // ""' values.yaml)"
      fi

      find . -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
      find "$SRC" -mindepth 1 -maxdepth 1 ! -name .git -exec cp -R {} . \;

      if [ -n "$tag" ] && [ -f values.yaml ]; then
        yq -i ".image.tag = \"$tag\"" values.yaml
      fi

      git_c add -A
      if git diff --cached --quiet; then
        echo "${each.key}: chart unchanged, nothing to push"
        exit 0
      fi
      git_c commit -qm "chart: sync ${each.value} from monorepo"
      git_c push -q origin HEAD:main
    EOT
  }
}
