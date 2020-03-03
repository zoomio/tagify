#!/bin/sh

DIST_DIR="_dist"
BINARY="tagify"
USER_BIN=$HOME/bin
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -z "$VERSION" ]; then
    VERSION="_$VERSION"
fi

file="${VERSION}/${BINARY}_${OS}_${VERSION}"

curl -O "https://github.com/zoomio/tagify/releases/download/${file}"
chmod +x ${file}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
  echo "created $USER_BIN directory, don't forget to add it to PATH environment variable"
fi

mv ${file} ${USER_BIN}/${BINARY}