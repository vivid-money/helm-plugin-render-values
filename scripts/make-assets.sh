#!/bin/bash
# set -x

VERSION="0.2.6"
mkdir -p assets

OS_LIST="darwin/arm64
darwin/amd64
linux/arm64
linux/amd64
windows/amd64
"

mkdir -p assets/plugin
cp plugin.yaml LICENSE assets/plugin/

for OS_ARCH in $(echo $OS_LIST)
do  
    OS=${OS_ARCH%/*}
    ARCH=${OS_ARCH#*/}
	env GOOS=${OS} GOARCH=${ARCH} go build -o assets/plugin/render-values .
	tar -czf assets/helm-plugin-render-values_${VERSION}_${OS}_${ARCH}.tar.gz -C assets/plugin .
done
