#!/bin/sh

TAG="$1"

if [ -z "$TAG" ]; then
    echo "TAG argument required"
    exit 1
fi

git tag "v$TAG"
git push origin "v$TAG"