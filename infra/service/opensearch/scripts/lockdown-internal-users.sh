#!/usr/bin/env bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: func ${FUNCNAME[0]:-main}" >&2' ERR
set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

NAMESPACE="${NAMESPACE:-opensearch}"
STATEFULSET="${STATEFULSET:-opensearch-cluster-master}"
NODE_POD="${NODE_POD:-opensearch-cluster-master-0}"
NODE_CONTAINER="${NODE_CONTAINER:-opensearch}"
HELPER_POD="${HELPER_POD:-opensearch-security-lockdown}"
ADMIN_SECRET="${ADMIN_SECRET:-opensearch-admin-credentials}"
KIBANASERVER_SECRET="${KIBANASERVER_SECRET:-opensearch-dashboards-kibanaserver}"
KIBANASERVER_PASSWORD_KEY="${KIBANASERVER_PASSWORD_KEY:-OPENSEARCH_KIBANASERVER_PASSWORD}"
REMOVED_USERS="${REMOVED_USERS:-logstash snapshotrestore kibanaro readall anomalyadmin}"
BACKUP_DIR="${BACKUP_DIR:-${PWD}/opensearch-security-backup-$(date -u +%Y%m%dT%H%M%SZ)}"
DRY_RUN="${DRY_RUN:-false}"

read -r -d '' PROBE_SCRIPT <<'EOF' || true
set -o errexit
set -o nounset
set -o pipefail
IFS= read -r encoded || true
test -n "${encoded}"
basic() {
  printf 'Authorization: Basic %s\n' "$(printf '%s:%s' "$1" "$2" | base64 -w0)"
}
status() {
  curl -s -o /dev/null -w '%{http_code}' -H @<(basic "$1" "$2") http://localhost:9200/_plugins/_security/authinfo
}
for user in ${REMOVED_USERS}; do
  echo "${user} $(status "${user}" "${user}")"
done
echo "kibanaserver-demo $(status kibanaserver kibanaserver)"
echo "kibanaserver-vault $(status kibanaserver "$(printf '%s' "${encoded}" | base64 -d)")"
source /vault/secrets/env-admin
reserved="$(curl -s -H @<(basic admin "${OPENSEARCH_INITIAL_ADMIN_PASSWORD}") http://localhost:9200/_plugins/_security/api/internalusers/kibanaserver |
  python3 -c 'import json, sys; print(str(json.load(sys.stdin)["kibanaserver"]["reserved"]).lower())')"
echo "kibanaserver-reserved ${reserved}"
EOF

read -r -d '' JOIN_SCRIPT <<'EOF' || true
set -o errexit
set -o nounset
set -o pipefail
config=/usr/share/opensearch/config
for _ in $(seq 1 60); do
  health="$(curl -sk -o /dev/null -w '%{http_code}' --cert "${config}/kirk.pem" --key "${config}/kirk-key.pem" https://localhost:9200/_plugins/_security/health || true)"
  nodes="$(curl -sk --cert "${config}/kirk.pem" --key "${config}/kirk-key.pem" 'https://localhost:9200/_cat/nodes?h=name' || true)"
  if [ "${health}" = "200" ] && awk -v name="${HELPER_POD}" '$1 == name { found = 1 } END { exit !found }' <<<"${nodes}"; then
    echo "Helper node joined the cluster and accepts the super admin certificate"
    exit 0
  fi
  sleep 5
done
echo "Helper node did not join the cluster in time" >&2
exit 1
EOF

read -r -d '' PREPARE_SCRIPT <<'EOF' || true
set -o errexit
set -o nounset
set -o pipefail
tools=/usr/share/opensearch/plugins/opensearch-security/tools
config=/usr/share/opensearch/config
work=/tmp/lockdown
admin_args=(-icl -nhnv -cacert "${config}/root-ca.pem" -cert "${config}/kirk.pem" -key "${config}/kirk-key.pem" -h localhost -p 9200)
rm -rf "${work}"
mkdir -p "${work}/backup"
if ! "${tools}/securityadmin.sh" -backup "${work}/backup" "${admin_args[@]}" </dev/null >"${work}/backup.log" 2>&1; then
  grep -vi 'hash' "${work}/backup.log" >&2
  exit 1
fi
test -s "${work}/backup/internal_users.yml"

block() {
  awk -v want="$1" '/^[^ #-]/ { key = $0; sub(/:.*/, "", key); gsub(/"/, "", key); current = key } current == want { print }' "$2"
}

LOCKDOWN_HASH="$("${tools}/hash.sh" -env KIBANASERVER_PASSWORD </dev/null | grep -E '^\$2[aby]\$' | tail -n 1)"
test -n "${LOCKDOWN_HASH}"
export LOCKDOWN_HASH

awk -v removed="${REMOVED_USERS}" '
  BEGIN { n = split(removed, names, " "); for (i = 1; i <= n; i++) skip[names[i]] = 1 }
  /^[^ #-]/ { key = $0; sub(/:.*/, "", key); gsub(/"/, "", key); current = key }
  current in skip { next }
  current == "kibanaserver" && /^  hash:/ { print "  hash: \"" ENVIRON["LOCKDOWN_HASH"] "\""; next }
  current == "kibanaserver" && /^  reserved:/ { print "  reserved: false"; next }
  { print }
' "${work}/backup/internal_users.yml" >"${work}/internal_users.yml"
unset LOCKDOWN_HASH

test -n "$(block _meta "${work}/internal_users.yml")"
test -n "$(block admin "${work}/internal_users.yml")"
[ "$(block admin "${work}/backup/internal_users.yml")" = "$(block admin "${work}/internal_users.yml")" ]
[ "$(block kibanaserver "${work}/internal_users.yml" | grep -c '^  hash: ')" = "1" ]
[ "$(block kibanaserver "${work}/internal_users.yml" | grep -c '^  reserved: false$')" = "1" ]
for user in ${REMOVED_USERS}; do
  test -z "$(block "${user}" "${work}/internal_users.yml")"
done
awk '/^[^ #-]/ { key = $0; sub(/:.*/, "", key); gsub(/"/, "", key); print "kept: " key }' "${work}/internal_users.yml"
EOF

read -r -d '' UPLOAD_SCRIPT <<'EOF' || true
set -o errexit
set -o nounset
set -o pipefail
tools=/usr/share/opensearch/plugins/opensearch-security/tools
config=/usr/share/opensearch/config
work=/tmp/lockdown
admin_args=(-icl -nhnv -cacert "${config}/root-ca.pem" -cert "${config}/kirk.pem" -key "${config}/kirk-key.pem" -h localhost -p 9200)
status=0
"${tools}/securityadmin.sh" -f "${work}/internal_users.yml" -t internalusers "${admin_args[@]}" </dev/null >"${work}/upload.log" 2>&1 || status=$?
grep -vi 'hash' "${work}/upload.log" || true
exit "${status}"
EOF

kube() {
  kubectl -n "${NAMESPACE}" "$@"
}

helper() {
  kube exec "${HELPER_POD}" -c opensearch -- "$@"
}

probe() {
  kube get secret "${KIBANASERVER_SECRET}" -o "jsonpath={.data.${KIBANASERVER_PASSWORD_KEY}}" |
    kube exec -i "${NODE_POD}" -c "${NODE_CONTAINER}" -- env REMOVED_USERS="${REMOVED_USERS}" bash -c "${PROBE_SCRIPT}"
}

locked_down() {
  local report="$1"
  local user
  for user in ${REMOVED_USERS}; do
    grep -qx "${user} 401" <<<"${report}" || return 1
  done
  grep -qx "kibanaserver-demo 401" <<<"${report}" &&
    grep -qx "kibanaserver-vault 200" <<<"${report}" &&
    grep -qx "kibanaserver-reserved false" <<<"${report}"
}

cleanup() {
  kube delete pod "${HELPER_POD}" --ignore-not-found --wait=false >/dev/null
}

if [ "$(kube get secret "${KIBANASERVER_SECRET}" -o "jsonpath={.data.${KIBANASERVER_PASSWORD_KEY}}" | wc -c)" -lt 2 ]; then
  echo "Secret ${NAMESPACE}/${KIBANASERVER_SECRET} has no ${KIBANASERVER_PASSWORD_KEY}; apply service/opensearch module.secret first" >&2
  exit 1
fi

report="$(probe)"
echo "${report}"
if locked_down "${report}"; then
  echo "Internal users already locked down"
  exit 0
fi

image="$(kube get statefulset "${STATEFULSET}" -o jsonpath='{.spec.template.spec.containers[0].image}')"
cluster_name="$(kube get statefulset "${STATEFULSET}" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="cluster.name")].value}')"
seed_hosts="$(kube get statefulset "${STATEFULSET}" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="discovery.seed_hosts")].value}')"

kube delete pod "${HELPER_POD}" --ignore-not-found --wait=true >/dev/null
trap cleanup EXIT

kube apply -f - <<MANIFEST
apiVersion: v1
kind: Pod
metadata:
  name: ${HELPER_POD}
  namespace: ${NAMESPACE}
  labels:
    app.kubernetes.io/name: opensearch-security-lockdown
spec:
  serviceAccountName: opensearch
  restartPolicy: Never
  securityContext:
    runAsUser: 1000
    fsGroup: 1000
  containers:
    - name: opensearch
      image: ${image}
      command: ["bash", "-c"]
      args:
        - |
          echo 'node.roles: []' >> config/opensearch.yml
          exec ./opensearch-docker-entrypoint.sh opensearch
      env:
        - name: cluster.name
          value: "${cluster_name}"
        - name: discovery.seed_hosts
          value: "${seed_hosts}"
        - name: node.name
          value: "${HELPER_POD}"
        - name: network.host
          value: "0.0.0.0"
        - name: OPENSEARCH_JAVA_OPTS
          value: "-Xms256m -Xmx256m"
        - name: DISABLE_PERFORMANCE_ANALYZER_AGENT_CLI
          value: "true"
        - name: OPENSEARCH_INITIAL_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              name: ${ADMIN_SECRET}
              key: initial_password
        - name: KIBANASERVER_PASSWORD
          valueFrom:
            secretKeyRef:
              name: ${KIBANASERVER_SECRET}
              key: ${KIBANASERVER_PASSWORD_KEY}
      resources:
        requests:
          cpu: 100m
          memory: 768Mi
        limits:
          memory: 1536Mi
      securityContext:
        runAsNonRoot: true
        capabilities:
          drop: ["ALL"]
MANIFEST

kube wait --for=condition=Ready "pod/${HELPER_POD}" --timeout=300s >/dev/null
helper env HELPER_POD="${HELPER_POD}" bash -c "${JOIN_SCRIPT}"
helper env REMOVED_USERS="${REMOVED_USERS}" bash -c "${PREPARE_SCRIPT}"

(
  umask 077
  mkdir -p "${BACKUP_DIR}"
  helper tar -C /tmp/lockdown/backup -cf - . | tar -C "${BACKUP_DIR}" -xf -
)
echo "Security configuration backup written to ${BACKUP_DIR} (contains password hashes; keep it private)"

if [ "${DRY_RUN}" = "true" ]; then
  echo "DRY_RUN=true: internal_users.yml rewritten and validated in ${HELPER_POD}, upload skipped"
  exit 0
fi

helper bash -c "${UPLOAD_SCRIPT}"

for _ in $(seq 1 12); do
  report="$(probe)"
  if locked_down "${report}"; then
    echo "${report}"
    echo "Internal users locked down; roll out OpenSearch Dashboards with the Vault kibanaserver password now"
    exit 0
  fi
  sleep 5
done
echo "${report}"
echo "Lockdown verification failed" >&2
exit 1
