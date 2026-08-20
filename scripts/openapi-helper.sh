#!/usr/bin/env bash
trap 'echo "${bash_source[0]}: line ${lineno}: status ${?}: user ${user}: func ${funcname[0]}"' err
set -o errexit
set -o errtrace
set -o xtrace

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

(cd "${PROJECT_ROOT}/services/admin" && make openapi)
(cd "${PROJECT_ROOT}/services/catalog" && make openapi)
(cd "${PROJECT_ROOT}/services/customer" && make openapi)
(cd "${PROJECT_ROOT}/services/delivery" && make openapi)
(cd "${PROJECT_ROOT}/services/identity" && make openapi)
(cd "${PROJECT_ROOT}/services/inventory" && make openapi)
(cd "${PROJECT_ROOT}/services/order" && make openapi)
(cd "${PROJECT_ROOT}/services/payment" && make openapi)
(cd "${PROJECT_ROOT}/services/notification" && make openapi)
(cd "${PROJECT_ROOT}/services/settlement" && make openapi)
(cd "${PROJECT_ROOT}/services/support" && make openapi)
