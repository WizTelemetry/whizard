#!/usr/bin/env bash

# 
# Copyright 2024 the Whizard Authors.
# 
# Licensed under Apache License, Version 2.0 with a few additional conditions.
# 
# You may obtain a copy of the License at
# 
#     https://github.com/WhizardTelemetry/whizard/blob/main/LICENSE
# 

# GoFmt apparently is changing @ head...

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::verify_go_version

cd "${KUBE_ROOT}"

find_files() {
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './.git' \
        -o -wholename './_output' \
        -o -wholename './_gopath' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/third_party/*' \
        -o -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}

find_files | xargs gofmt -s -w