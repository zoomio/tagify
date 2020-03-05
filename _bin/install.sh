#!/bin/sh

BINARY="tagify"
USER_BIN=$HOME/bin
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -z "$VERSION" ]; then
    VERSION="_$VERSION"
fi

link=$(curl -s https://api.github.com/repos/jgm/pandoc/releases/latest | grep "browser_download_url.*deb" | cut -d : -f 2,3 | tr -d \")

curl -o ${BINARY} ${link}
chmod +x ${BINARY}

if [ ! -d "$USER_BIN" ]; then
  mkdir -p ${USER_BIN}
  echo "created $USER_BIN directory, don't forget to add it to PATH environment variable"
fi

mv ${BINARY} ${USER_BIN}/${BINARY}

echo "Installation is done."