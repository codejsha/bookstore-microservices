#!/usr/bin/env bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace
set -o xtrace

ISTIO_VERSION="1.24.3"

declare -A DASHBOARDS=(
  [7639]="mesh"
  [7636]="service"
  [7630]="workload"
  [11829]="performance"
  [7645]="controlplane"
  [13277]="wasm"
)

mkdir -p dashboards
for ID in "${!DASHBOARDS[@]}"; do
  NAME="${DASHBOARDS[$ID]}"
  REVISION=$(curl -s "https://grafana.com/api/dashboards/${ID}/revisions" \
    | jq -r ".items[] | select(.description | contains(\"${ISTIO_VERSION}\")) | .revision" \
    | tail -n 1)

  if [ ! -f "dashboards/istio-${NAME}.json" ]; then
    curl -s "https://grafana.com/api/dashboards/${ID}/revisions/${REVISION}/download" \
      -o "dashboards/istio-${NAME}.json"
  fi
done
