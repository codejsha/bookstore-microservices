#!/usr/bin/env bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace
set -o xtrace

ISTIO_VERSION="1.29.1"

declare -A DASHBOARDS=(
  [7639]="mesh"
  [7636]="service"
  [7630]="workload"
  [11829]="performance"
  [7645]="controlplane"
  [13277]="wasm"
  [21306]="ztunnel"
)

ISTIO_DIR="dashboards/Infra - Mesh"
mkdir -p "${ISTIO_DIR}"
for ID in "${!DASHBOARDS[@]}"; do
  NAME="${DASHBOARDS[$ID]}"
  REVISION=$(curl -s "https://grafana.com/api/dashboards/${ID}/revisions" \
    | jq -r ".items[] | select(.description | contains(\"${ISTIO_VERSION}\")) | .revision" \
    | tail -n 1)

  if [ ! -f "${ISTIO_DIR}/istio-${NAME}.json" ]; then
    curl -s "https://grafana.com/api/dashboards/${ID}/revisions/${REVISION}/download" \
      -o "${ISTIO_DIR}/istio-${NAME}.json"
  fi
done

# ----- Bookstore infrastructure dashboards from grafana.com -----
declare -A INFRA_DASHBOARDS=(
  [20417]="Infra - Data/postgres-cnpg"        # CloudNativePG official
  [14057]="Infra - Data/mysql-innodb"         # mysqld_exporter
  [11835]="Infra - Data/valkey"               # redis_exporter
  [18276]="Infra - Streaming/kafka-cluster"   # Kafka cluster (JMX exporter)
  [14584]="Infra - GitOps/argocd"             # ArgoCD official
  [13192]="Infra - GitOps/gitea"              # Gitea
  [17878]="Infra - Platform/keycloak"         # Keycloak metrics
  [12904]="Infra - Platform/vault"            # Hashicorp Vault telemetry
  [11001]="Infra - Platform/cert-manager"     # cert-manager
  [14075]="Infra - GitOps/harbor"             # Harbor native metrics
  [10423]="Infra - Data/seaweedfs"            # SeaweedFS
)

fetch_clean() {
  curl -s "$1" | tr -d '\000-\037'
}

for ID in "${!INFRA_DASHBOARDS[@]}"; do
  REL_PATH="${INFRA_DASHBOARDS[$ID]}"
  OUT="dashboards/${REL_PATH}.json"
  mkdir -p "$(dirname "${OUT}")"
  if [ -f "${OUT}" ]; then
    continue
  fi

  REVISION=$(fetch_clean "https://grafana.com/api/dashboards/${ID}/revisions" \
    | jq -r '(.items // []) | sort_by(.revision) | (.[-1].revision // empty)')

  if [ -z "${REVISION}" ]; then
    echo "WARN: no revision for grafana.com dashboard ${ID}" >&2
    continue
  fi

  fetch_clean "https://grafana.com/api/dashboards/${ID}/revisions/${REVISION}/download" \
    | jq '
        (.. | objects | select(has("datasource")) | .datasource) |=
          (if (type == "object") then (. + {uid: "prometheus", type: "prometheus"})
           else "prometheus" end)
        | if has("__inputs") then
            .__inputs |= map(
              if .name == "DS_PROMETHEUS" then . + {pluginId: "prometheus", type: "datasource"} else . end
            )
          else . end
        | del(.__elements)
      ' > "${OUT}"
done

# ----- Bookstore infrastructure dashboards from upstream GitHub -----
declare -A GITHUB_DASHBOARDS=(
  ["Infra - Streaming/kafka-connect"]="https://raw.githubusercontent.com/strimzi/strimzi-kafka-operator/main/examples/metrics/grafana-dashboards/strimzi-kafka-connect.json"
  ["Infra - Streaming/kafka-exporter"]="https://raw.githubusercontent.com/strimzi/strimzi-kafka-operator/main/examples/metrics/grafana-dashboards/strimzi-kafka-exporter.json"
  ["Infra - GitOps/argo-rollouts"]="https://raw.githubusercontent.com/argoproj/argo-rollouts/master/examples/dashboard.json"
  ["Infra - Platform/temporal"]="https://raw.githubusercontent.com/temporalio/dashboards/master/server/server-general.json"
)

for REL_PATH in "${!GITHUB_DASHBOARDS[@]}"; do
  URL="${GITHUB_DASHBOARDS[$REL_PATH]}"
  OUT="dashboards/${REL_PATH}.json"
  mkdir -p "$(dirname "${OUT}")"
  if [ -f "${OUT}" ]; then
    continue
  fi

  fetch_clean "${URL}" \
    | jq '
        (.. | objects | select(has("datasource")) | .datasource) |=
          (if (type == "object") then (. + {uid: "prometheus", type: "prometheus"})
           else "prometheus" end)
        | if has("__inputs") then
            .__inputs |= map(
              if .name == "DS_PROMETHEUS" then . + {pluginId: "prometheus", type: "datasource"} else . end
            )
          else . end
        | del(.__elements)
      ' > "${OUT}"
done
