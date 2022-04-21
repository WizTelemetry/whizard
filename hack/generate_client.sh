#!/bin/bash

set -e

GV="$1"

rm -rf ./pkg/client/clientset ./pkg/client/informers ./pkg/client/listers
./hack/generate_group.sh "client,lister,informer" github.com/kubesphere/paodin-monitoring/pkg/client github.com/kubesphere/paodin-monitoring/pkg/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv github.com/kubesphere/paodin-monitoring/pkg/client/* ./pkg/client/
rm -rf ./github.com
