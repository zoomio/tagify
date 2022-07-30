#!/bin/sh

DIST_DIR="_dist"
BINARY="tagify"
OS="$1"
VERSION="$2"
ARCH="$3"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ -z "$VERSION" ]; then
    VERSION="tip"
fi

if [ -z "$ARCH" ]; then
    ARCH="amd64"
fi

if [ ! -d "$DIST_DIR" ]; then
  mkdir -p $DIST_DIR
fi

env GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-X main.version=${VERSION}" -v -o ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION} cmd/cli/cli.go