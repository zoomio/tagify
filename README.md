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

Want to see it in action? Visit [Tagify Playground](https://www.zoomio.org/tagify).

Example, "tagify" this repository (with the limit of 5 tags):
```bash
tagify -s https://github.com/zoomio/tagify -l 5
```

In a code (see [cmd/cli/cli.go](https://raw.githubusercontent.com/zoomio/tagify/master/cmd/cli/cli.go)).

Use `-no-stop` flag to disable filtering out of the [stop-words](https://github.com/zoomio/stopwords).

## Extensions (Beta)

Since `v0.50.0` Tagify has added support for extensions. See `extension/extension.go` and its usages and implementations in `processor/html/extension.go`. You can see an example at `processor/html/extension_test.go`.

## Installation

### Binary

Get the latest [release](https://github.com/zoomio/tagify/releases/latest) by running this command in your shell:

__For MacOS:__
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" -o darwin
```

__For MacOS (arm64):__
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" -o darwin arm64
```

__For Linux:__
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" -o linux
```

__For Windows:__
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/tagify/master/_bin/install.sh)" -o windows
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

[![Buy Me A Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/smeshkov)
