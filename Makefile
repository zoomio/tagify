.PHONY: deps clean build

TAG=0.60.2
BINARY=tagify
DIST_DIR=_dist
OS=darwin
ARCH=arm64
VERSION=tip
USER_BIN=${HOME}/bin
DATE=`date +%m-%d-%Y-%H-%M-%S`

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

bench:
	go test github.com/zoomio/tagify/processor/html -bench=. -run=ParseHTML -count=5 | tee _dist/parse_html_${DATE}.txt

tag:
	./_bin/tag.sh ${TAG}

# install:
# 	./_bin/install.sh ${OS} ${ARCH}

install: build
	chmod +x ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION}
	mv ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION} ${USER_BIN}/${BINARY}