#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

UNSEAL_KEY="$(jq -r '.unseal_keys_b64[0]' ./cluster-keys.json)"
kubectl exec -n vault vault-0 -- vault operator unseal ${UNSEAL_KEY}
kubectl exec -n vault vault-1 -- vault operator unseal ${UNSEAL_KEY}
kubectl exec -n vault vault-2 -- vault operator unseal ${UNSEAL_KEY}
