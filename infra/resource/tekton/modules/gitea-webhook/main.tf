terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
    }
  }
}

resource "null_resource" "webhook" {
  for_each = var.webhooks

  triggers = {
    hooks_api = "${var.gitea_api_url}/repos/${var.org_name}/${each.key}${var.source_repo_suffix}/hooks"
    mode      = "absent"
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    environment = {
      GITEA_USER = var.gitea_username
      GITEA_PASS = var.gitea_password
      HOOKS_API  = self.triggers.hooks_api
    }
    command = <<-EOT
      set -euo pipefail

      api() {
        local method="$1" url="$2" expect="$3"
        local out code payload
        out=$(curl -s -w '\n%%{http_code}' -u "$GITEA_USER:$GITEA_PASS" -X "$method" "$url")
        code=$${out##*$'\n'}
        payload=$${out%$'\n'*}
        case ",$expect," in
          *",$code,"*) printf '%s' "$payload"; return 0 ;;
        esac
        echo "gitea-webhook: $method $url -> HTTP $code: $payload" >&2
        return 1
      }

      api GET "$HOOKS_API?limit=50" 200 \
        | python3 -c 'import sys,json
from urllib.parse import urlparse
try: hooks=json.load(sys.stdin)
except Exception: hooks=[]
for h in hooks:
    host=urlparse(h.get("config",{}).get("url","")).hostname or ""
    if host.startswith("el-"): print(h["id"])' \
        | while read -r id; do
            [ -n "$id" ] && api DELETE "$HOOKS_API/$id" 204,200 >/dev/null
          done
    EOT
  }
}
