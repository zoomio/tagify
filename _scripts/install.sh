#!/bin/sh

BINARY="tagify"
USER_BIN=$HOME/bin

if [ ! -f "$BINARY" ]; then
    echo "$BINARY not found"
    exit 1
fi

chmod +x ${BINARY}
mv ${BINARY} ${USER_BIN}/${BINARY}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
fi
