# Tagify

[![Go Report Card](https://goreportcard.com/badge/github.com/zoomio/tagify)](https://goreportcard.com/report/github.com/zoomio/tagify)
[![Coverage](https://codecov.io/gh/zoomio/tagify/branch/master/graph/badge.svg)](https://codecov.io/gh/zoomio/tagify)
[![GoDoc](https://godoc.org/github.com/zoomio/tagify?status.svg)](https://godoc.org/github.com/zoomio/tagify)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Tagify can recieve STDIN, file or HTTP address as an input and return a list of most popular words ordered by popularity as an output.

More info about what is "Tagify" and the reasons behind it can be found [here](https://zoomio.org/blog/post/tags_as_a_service-5712840111423488).

Supported formats:
- Plain text
- HTML
- Markdown

Supported languages:
- English
- Russian
- Chinese
- Hindi
- Hebrew
- Spanish
- Arabic
- Japanese
- German
- French
- Korean

Visit Tagify playground [here](https://www.zoomio.org/tagify):

![ZoomIO Tagify](https://storage.googleapis.com/www.zoomio.org/ZoomIO_tagify.png)

Example, "tagify" this repository (with the limit of 5 tags):
```bash
$ tagify -s=https://github.com/zoomio/tagify -l=5
source html plain supports tags
```

In a code (see [cmd/cli/cli.go](https://raw.githubusercontent.com/zoomio/tagify/master/cmd/cli/cli.go)).

Use `-no-stop` flag to disable filtering out of the [stop-words](https://github.com/zoomio/stopwords).

## Installation

### Binary

Get the latest [release](https://github.com/zoomio/tagify/releases/latest) by running this command in your shell:

For MacOS:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" darwin
```

For MacOS (arm64):
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" darwin arm64
```

For Linux:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" linux
```

For Windows:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" windows
```

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
