#!/bin/bash

VERSION="0.1.1"
mkdir -p assets

# MACOS
mkdir assets/helm-plugin-render-values_${VERSION}_darwin_amd64
cp plugin.yaml LICENSE assets/helm-plugin-render-values_${VERSION}_darwin_amd64/
env GOOS=darwin go build -o assets/helm-plugin-render-values_${VERSION}_darwin_amd64/render-values .

# LINUX
mkdir assets/helm-plugin-render-values_${VERSION}_linux_amd64
cp plugin.yaml LICENSE assets/helm-plugin-render-values_${VERSION}_linux_amd64/
env GOOS=linux go build -o assets/helm-plugin-render-values_${VERSION}_linux_amd64/render-values .

# LINUX
mkdir assets/helm-plugin-render-values_${VERSION}_windows_amd64
cp plugin.yaml LICENSE assets/helm-plugin-render-values_${VERSION}_windows_amd64/
env GOOS=windows go build -o assets/helm-plugin-render-values_${VERSION}_windows_amd64/render-values .

cd assets
tar -czf helm-plugin-render-values_${VERSION}_linux_amd64.tar.gz helm-plugin-render-values_${VERSION}_linux_amd64
tar -czf helm-plugin-render-values_${VERSION}_darwin_amd64.tar.gz helm-plugin-render-values_${VERSION}_darwin_amd64
tar -czf helm-plugin-render-values_${VERSION}_windows_amd64.tar.gz helm-plugin-render-values_${VERSION}_windows_amd64
