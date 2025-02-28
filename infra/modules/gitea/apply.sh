#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

terraform init -upgrade
terraform apply -auto-approve -target module.helm
terraform apply -auto-approve -target module.istio
terraform apply -auto-approve -target module.ssh
terraform apply -auto-approve -target module.user
terraform apply -auto-approve -target module.organization
terraform apply -auto-approve -target module.repository
terraform apply -auto-approve -target module.team
