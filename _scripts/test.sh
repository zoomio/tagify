#!/bin/sh

go get -u github.com/alecthomas/gometalinter
gometalinter --install
go install cmd/cli/cli.go
gometalinter --fast --vendor ./...

go test -v ./...