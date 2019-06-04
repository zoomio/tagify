# Tagify

[![Build Status](https://travis-ci.org/zoomio/tagify.svg?branch=master)](https://travis-ci.org/zoomio/tagify)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoomio/tagify)](https://goreportcard.com/report/github.com/zoomio/tagify)

Visit playground [here](https://www.zoomio.org/tagify).

Gets STDIN, file or HTTP address as an input and returns ordered list of most frequent words as an output. Does not support Single Page Apps (browser rendered) at the moment - Work-in-Progress.

Example, get 10 most frequent words from StackOverflow main page:
```bash
$ tagify -s=https://github.com/zoomio/tagify -l=3
frequent address apache
```

In code (see [cmd/cli/cli.go](https://raw.githubusercontent.com/zoomio/tagify/master/cmd/cli/cli.go)).

Use `-no-stop` flag to disable stop-words filtering ([processor/stopwords.go](https://raw.githubusercontent.com/zoomio/tagify/master/processor/stopwords.go)).

## Installation

### Binary

* download latest release for corrseponding OS (Darwin or Linux) from [Releases](https://github.com/zoomio/tagify/releases/latest)
* make binary executable: `chmod +x <binary>`
* put executable binary under your bin directory, e.g. (assuming `~/bin` is in your PATH): `mv <binary> $HOME/bin/<binary>`.

### Go dependency

```bash
go get -u github.com/zoomio/tagify/...
```

## Changelog

See [CHANGELOG.md](https://raw.githubusercontent.com/zoomio/tagify/master/CHANGELOG.md)

## Contributing

See [CONTRIBUTING.md](https://raw.githubusercontent.com/zoomio/tagify/master/CONTRIBUTING.md)

## License

Released under the [Apache License 2.0](https://raw.githubusercontent.com/zoomio/tagify/master/LICENSE).