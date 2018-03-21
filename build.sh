#!/bin/sh

OS="$1"

if [ -z "$OS" ]; then
    OS="darwin"
fi

go test ./...
env GOOS=${OS} GOARCH=amd64 go build -o tagify cmd/cli/cli.go