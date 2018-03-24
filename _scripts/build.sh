#!/bin/sh

OS="$1"

if [ -z "$OS" ]; then
    OS="darwin"
fi

go test -v ./...

# use packr, to include files in binary"
go get -u github.com/gobuffalo/packr/...

env GOOS=${OS} GOARCH=amd64 packr build -o tagify_${OS} cmd/cli/cli.go