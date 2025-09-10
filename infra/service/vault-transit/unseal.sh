#!/usr/bin/env bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

UNSEAL_KEY="$(jq -r '.unseal_keys_b64[0]' ./unsealer-keys.json)"
kubectl exec -n vault-transit vault-0 -- vault operator unseal "${UNSEAL_KEY}"
