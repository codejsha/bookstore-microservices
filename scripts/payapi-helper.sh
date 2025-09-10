#!/usr/bin/env bash
trap 'echo "${bash_source[0]}: line ${lineno}: status ${?}: user ${user}: func ${funcname[0]}"' err
set -o errexit
set -o errtrace
set -o xtrace

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

(cd "${PROJECT_ROOT}/services/payment" && make payapi)
