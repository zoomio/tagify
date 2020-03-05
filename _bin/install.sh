#!/bin/sh

BINARY="tagify"
USER_BIN=$HOME/bin
OS="$1"

if [ -z "$OS" ]; then
    OS="darwin"
fi

link=$(curl -s https://api.github.com/repos/zoomio/tagify/releases/latest | grep "browser_download_url.*tagify_${OS}" | cut -d : -f 2,3 | tr -d \")

curl -o ${BINARY} ${link}
chmod +x ${BINARY}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
  echo "created $USER_BIN directory, don't forget to add it to PATH environment variable"
fi

mv ${BINARY} ${USER_BIN}/${BINARY}

echo "Installation is done."