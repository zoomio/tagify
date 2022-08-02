.PHONY: deps clean build

TAG=0.57.0
BINARY=tagify
DIST_DIR=_dist
OS=darwin
ARCH=arm64
VERSION=tip
USER_BIN=${HOME}/bin

deps:
	go get -u ./...

clean: 
	rm -rf _dist/*
	
build:
	./_bin/build.sh ${OS} ${VERSION} ${ARCH}

test:
	./_bin/test.sh

run:
	./_dist/tagify_darwin -s=https://zoomio.org/blog/post/mock_server-5632006343884800

profile:
	./_dist/tagify_darwin -s=https://zoomio.org/blog/post/mock_server-5632006343884800 -cpuprofile=_dist/tagify_darwin.prof
	go tool pprof _dist/tagify_darwin _dist/tagify_darwin.prof

tag:
	./_bin/tag.sh ${TAG}

install:
	./_bin/install.sh ${OS} ${ARCH}

install_local: build
	chmod +x ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION}
	mv ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION} ${USER_BIN}/${BINARY}