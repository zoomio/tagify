# Contributing

## Guidelines for pull requests

- Write tests for any changes.
- Separate unrelated changes into multiple pull requests.
- For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.
- Ensure the build is green before you open your PR. The Pipelines build won't run by default on a remote branch, so enable Pipelines.

## Build

* [Go](https://golang.org/dl/)
* To build binary run `./_bin/build.sh` in shell, it will produce `tagify` binary.
* To install use `./_bin/install.sh`, it will put `tagify` binary under `~/bin` directory assuming it is in your `PATH`.

## Release

* All notable changes comming with the new version should be documented in [CHANGELOG.md](https://raw.githubusercontent.com/zoomio/tagify/master/CHANGELOG.md).
* Run tests with `./_bin/test.sh`, make sure everything is passing.
* Tag, push and trigger new binary release on GitHub `./_bin/tag.sh [version]`.
* To perform Brew release, use `./_bin/brew_release.sh v[version]`, then submit a PR to Homebrew repo with the file from `./_templates/tagify.rb`.