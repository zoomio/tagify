# Tagify

[![Build Status](https://travis-ci.org/zoomio/tagify.svg?branch=master)](https://travis-ci.org/zoomio/tagify)

Gets STDIN, file or HTTP address as an input and returns ordered list of most frequent tags as an output.

Example, get 10 most frequent words from StackOverflow main page:
```bash
$ tagify -s=https://stackoverflow.com -l=10
application using page add file server run ionic local error
```

## Installation

### Binary

* download latest release for corrseponding OS (Darwin or Linux) from [Releases](https://github.com/zoomio/tagify/releases/latest)
* make binary executable: `chmod +x <binary>`
* put executable binary under your bin directory, e.g. (assuming `~/bin` is in your PATH): `mv <binary> $HOME/bin/<binary>`.

### Go dependency

```bash
go get -u github.com/zoomio/tagify/...
```

## Development

* [Go](https://golang.org/dl/)
* [Dep](https://golang.github.io/dep/docs/installation.html)

## License

Released under the [Apache License 2.0](https://raw.githubusercontent.com/zoomio/tagify/master/LICENSE).