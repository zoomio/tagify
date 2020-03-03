#!/bin/sh

DIST_DIR="_dist"
BINARY="tagify"
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ -z "$VERSION" ]; then
    VERSION="tip"
fi

if [ ! -d "$DIST_DIR" ]; then
  mkdir -p $DIST_DIR
fi

env GOOS=${OS} GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -v -o ${DIST_DIR}/${BINARY}_${OS}_${VERSION} cmd/cli/cli.go