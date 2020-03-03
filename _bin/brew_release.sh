#!/bin/sh

version="$1"

if [ -z "$version" ]; then
    echo "version is required, e.g. `${0} [version]`"
    exit 2
fi

tagifyCode="_dist/tagify_code.tar.gz"

curl -o ${tagifyCode} "https://codeload.github.com/zoomio/tagify/tar.gz/$version"

export BREW_BIN=$(brew --prefix)/bin
export VERSION=$version
export SHA=$(openssl sha256 < ${tagifyCode})

envsubst < _templates/tagify.template.rb > _dist/tagify.rb

echo "ok"