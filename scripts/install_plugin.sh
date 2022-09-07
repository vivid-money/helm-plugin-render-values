#!/bin/sh -e
set -xe

# Copied w/ love from the excellent hypnoglow/helm-s3

if [ -n "${HELM_PUSH_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
echo "Downloading and installing helm-push v${version} ..."

arch=""
url=""
if [ "$(arch)" = "aarch64" ]; then
    arch="arm64"
elif [ "$(arch)" = "arm" ] ; then
    arch="arm64"
else
    arch="amd64"
fi


if [ "$(uname)" = "Darwin" ]; then
    url="https://github.com/vivid-money/helm-plugin-render-values/releases/download/v${version}/helm-plugin-render-values_${version}_darwin_${arch}.tar.gz"
elif [ "$(uname)" = "Linux" ] ; then
    url="https://github.com/vivid-money/helm-plugin-render-values/releases/download/v${version}/helm-plugin-render-values_${version}_linux_${arch}.tar.gz"
else
    url="https://github.com/vivid-money/helm-plugin-render-values/releases/download/v${version}/helm-plugin-render-values_${version}_windows_${arch}.tar.gz"
fi

echo $url

mkdir -p "releases/v${version}"

# Download with curl if possible.
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/v${version}.tar.gz"
else
    wget -q "${url}" -O "releases/v${version}.tar.gz"
fi
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"
mv "releases/v${version}/render-values" "render-values" || \
    mv "releases/v${version}/render-values.exe" "render-values"
