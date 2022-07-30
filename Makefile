.PHONY: deps clean build

TAG=0.56.0

deps:
	go get -u ./...

clean: 
	rm -rf _dist/*
	
build:
	./_bin/build.sh

test:
	./_bin/test.sh

run:
	./_dist/tagify_darwin -s=https://zoomio.org/blog/post/mock_server-5632006343884800

profile:
	./_dist/tagify_darwin -s=https://zoomio.org/blog/post/mock_server-5632006343884800 -cpuprofile=_dist/tagify_darwin.prof
	go tool pprof _dist/tagify_darwin _dist/tagify_darwin.prof

tag:
	./_bin/tag.sh ${TAG}
