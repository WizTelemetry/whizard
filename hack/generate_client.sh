#!/bin/bash

set -e

GV="$1"

rm -rf ./pkg/client/clientset ./pkg/client/informers ./pkg/client/listers
./hack/generate_group.sh "client,lister,informer" github.com/kubesphere/whizard/pkg/client github.com/kubesphere/whizard/pkg/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv github.com/kubesphere/whizard/pkg/client/* ./pkg/client/
rm -rf ./github.com
