#!/bin/bash
trap 'echo "${bash_source[0]}: line ${lineno}: status ${?}: user ${user}: func ${funcname[0]}"' err
set -o errexit
set -o errtrace
set -o xtrace

cd ../../services/catalog && make openapi
cd ../../services/customer && make openapi
cd ../../services/identity && make openapi
cd ../../services/inventory && make openapi
cd ../../services/order && make openapi
cd ../../services/payment && make openapi
