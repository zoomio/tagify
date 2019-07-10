#!/bin/sh

# formater
gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

# linter
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
golangci-lint run

# tests & coverage
go test -coverprofile=coverage.out -v ./...
go tool cover -func=coverage.out

# clean after self
go mod tidy