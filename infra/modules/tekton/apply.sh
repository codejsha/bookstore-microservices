#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace
set -o xtrace

/bin/cp -f ../vault/int-ca.crt .
ROOT_TOKEN=$(jq -r '.root_token' ../vault/cluster-keys.json)
perl -pi -e "s/^vault_token.*/vault_token = \"${ROOT_TOKEN}\"/" terraform.tfvars

terraform init -upgrade
terraform apply -auto-approve -target module.component
terraform apply -auto-approve -target module.istio

#terraform apply -auto-approve -target module.service_account
#terraform apply -auto-approve -target module.trigger_binding
#terraform apply -auto-approve -target module.golnag_template
