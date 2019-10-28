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

PATH=${DIST_DIR}/${BINARY}_${OS}${VERSION}

chmod +x ${PATH}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
  echo "created $USER_BIN directory, don't forget to add it to PATH environment variable"
fi

mv ${PATH} ${USER_BIN}/${BINARY}