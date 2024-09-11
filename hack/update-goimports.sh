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

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"
source "${KUBE_ROOT}/hack/lib/util.sh"

kube::golang::verify_go_version

# Ensure that we find the binaries we build before anything else.
export GOBIN="${KUBE_OUTPUT_BINPATH}"
PATH="${GOBIN}:${PATH}"

# Explicitly opt into go modules, even though we're inside a GOPATH directory
export GO111MODULE=on

if ! command -v goimports ; then
# Install goimports
  echo 'installing goimports'
  pushd "${KUBE_ROOT}/hack/tools" >/dev/null
    GO111MODULE=auto go install -mod=mod golang.org/x/tools/cmd/goimports@v0.7.0
  popd >/dev/null
fi

cd "${KUBE_ROOT}" || exit 1

IFS=$'\n' read -r -d '' -a files < <( find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./pkg/client/*" -not -name "zz_generated.deepcopy.go" && printf '\0' )

"goimports" -w -local github.com/WhizardTelemetry/whizard "${files[@]}"