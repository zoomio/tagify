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

# use packr, to include files in binary
go get -u github.com/gobuffalo/packr/...

env GOOS=${OS} GOARCH=amd64 packr build -o ${BINARY}_${OS}${VERSION} cmd/cli/cli.go