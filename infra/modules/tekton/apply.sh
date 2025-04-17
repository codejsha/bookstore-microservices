#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

terraform init -upgrade
terraform apply -auto-approve -target module.component
terraform apply -auto-approve -target module.istio
terraform apply -auto-approve -target module.service_account
terraform apply -auto-approve -target module.trigger_binding
terraform apply -auto-approve -target module.golnag_template
