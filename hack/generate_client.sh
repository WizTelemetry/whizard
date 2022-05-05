#!/bin/bash

set -e

GV="$1"

rm -rf ./pkg/client/clientset ./pkg/client/informers ./pkg/client/listers
./hack/generate_group.sh "client,lister,informer" github.com/kubesphere/paodin/pkg/client github.com/kubesphere/paodin/pkg/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv github.com/kubesphere/paodin/pkg/client/* ./pkg/client/
rm -rf ./github.com
