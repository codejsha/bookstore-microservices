#!/usr/bin/env bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

VAULT_ADDR="${VAULT_ADDR:-http://vault.example.com}"
VAULT_PUBLISHER_PATH="${VAULT_PUBLISHER_PATH:-kv/nexus/publisher}"
NEXUS_REGISTRY="${NEXUS_REGISTRY:-nexus.example.com/repository/npm-hosted}"
NEXUS_CA_FILE="${NEXUS_CA_FILE:-}"
export VAULT_ADDR

die() {
  echo "ERROR: ${1}" >&2
  exit 1
}

[ "$#" -gt 0 ] || die "usage: $(basename "$0") <package-dir> [<package-dir>...]"

if [ -z "${VAULT_TOKEN:-}" ]; then
  [ -r "${HOME}/.vault/token" ] || die "no VAULT_TOKEN and ~/.vault/token is unreadable"
  VAULT_TOKEN="$(cat "${HOME}/.vault/token")"
fi
export VAULT_TOKEN

user="$(vault kv get -field=username "${VAULT_PUBLISHER_PATH}")" &&
  pass="$(vault kv get -field=password "${VAULT_PUBLISHER_PATH}")" ||
  die "cannot read ${VAULT_PUBLISHER_PATH} from Vault (check VAULT_ADDR / ~/.vault/token)"

npmrc="$(mktemp)"
ca="$(mktemp)"
trap 'rm -f "${npmrc}" "${ca}"' EXIT

printf '//%s/:_auth=%s\nalways-auth=true\n' \
  "${NEXUS_REGISTRY}" \
  "$(printf '%s:%s' "${user}" "${pass}" | base64 | tr -d '\n')" > "${npmrc}"

if [ -n "${NEXUS_CA_FILE}" ]; then
  cat "${NEXUS_CA_FILE}" > "${ca}" 2>/dev/null || true
else
  vault read -field=certificate pki_int/cert/ca_chain > "${ca}" 2>/dev/null || true
fi
[ -s "${ca}" ] || die "no CA cert (set NEXUS_CA_FILE=/path/to/ca.pem or ensure Vault pki_int is reachable)"

for dir in "$@"; do
  name="$(node -p "require('./${dir}/package.json').name" 2>/dev/null || echo "${dir}")"
  echo "==> publishing ${name}"
  (
    cd "${dir}" &&
      NPM_CONFIG_USERCONFIG="${npmrc}" NODE_EXTRA_CA_CERTS="${ca}" pnpm publish --no-git-checks
  )
done
