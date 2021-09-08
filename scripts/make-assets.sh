#!/bin/bash
set -x

VERSION="v0.1.2"
mkdir -p assets

OS_LIST="darwin
linux
windows
"

mkdir assets/plugin
cp plugin.yaml LICENSE assets/plugin/

for OS in $(echo $OS_LIST)
do
	env GOOS=${OS} go build -o assets/plugin/render-values .
	tar -czf assets/helm-plugin-render-values_${VERSION}_${OS}_amd64.tar.gz -C assets/plugin .
done
