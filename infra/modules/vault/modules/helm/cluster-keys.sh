#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

init_status=$(kubectl exec -n vault vault-0 -- vault status -format=json | jq -r '.initialized')
if [ "${init_status}" == "false" ]; then
    kubectl exec -n vault vault-0 -- vault operator init \
        -key-shares=1 \
        -key-threshold=1 \
        -format=json \
        > cluster-keys.json
fi

echo "{\"vault_root_token\": \"$(cat cluster-keys.json | jq -r .root_token)\", " \
    "\"vault_unseal_key\": \"$(cat cluster-keys.json | jq -r .unseal_keys_b64[0])\"}"
