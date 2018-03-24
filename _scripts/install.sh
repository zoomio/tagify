#!/bin/sh

BINARY="tagify"
USER_BIN=$HOME/bin
OS="$1"

if [ -z "$OS" ]; then
    OS="darwin"
fi

chmod +x ${BINARY}_${OS}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
fi

mv ${BINARY}_${OS} ${USER_BIN}/${BINARY}