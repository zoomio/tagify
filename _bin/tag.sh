#!/bin/sh

TAG="$1"

if [ -z "$TAG" ]; then
    echo "TAG is required, e.g. `${0} [tag]`"
    exit 1
fi

git tag "v$TAG"
git push origin "v$TAG"