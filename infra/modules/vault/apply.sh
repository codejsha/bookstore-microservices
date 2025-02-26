#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

terraform init -upgrade
terraform apply -auto-approve -target module.cert
terraform apply -auto-approve -target module.helm
terraform apply -auto-approve -target module.istio
terraform apply -auto-approve -target module.pki
terraform apply -auto-approve -target module.kv
terraform apply -auto-approve -target module.kube
