#!/bin/sh

gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

# lint code
go get -u github.com/alecthomas/gometalinter
gometalinter --install
go install cmd/cli/cli.go
gometalinter --fast --vendor ./...

# test code
go test -v ./...