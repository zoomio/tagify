#!/bin/sh

BINARY="tagify"
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -z "$VERSION" ]; then
    VERSION="_$VERSION"
fi

env GOOS=${OS} GOARCH=amd64 go build -o ${BINARY}_${OS}${VERSION} cmd/cli/cli.go