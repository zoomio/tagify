# Changelog

## v0.45.0

- HTML: prioritize longer page titles over the shorter ones;
- bumped Go version to 1.16;
- released with GitHub Actions.

## 0.42.1
 - added support for Chinese, Hindi, Spanish, Arabic, Japanese, German, Hebrew, French and Korean tags.

## 0.41.2
 - a milestone release (hence the jump in the version number), part of the ["Faster Stronger Better"](https://zoomio.org/blog/post/faster_stronger_better-5708658021236736) initiative;
 - added support for Markdown content type;
 - improved performance and accuracy of HTML tagifier;
 - added **experiemntal** `-content` option - it allows to target "content" only tags (such as headings and paragraphs);
 - added **experiemntal** `-site` option - it allows to Tagify full site;
 - added `Result#ForEach` for easier iteration through the tags;
 - added `-version` option to show version of Tagify.

## 0.35.0
 - BREAKING CHANGE (with `v0.33.0`): renamed `Result.Meta.DocVersion` to `Result.Meta.DocHash`.

## 0.34.0
 - unified `#GetTagsFromString` with `#Run` so it is now a single API call - `#Run`;
 - fixed logic with `Query` option when it was wrongly setting `ContentType` to `Text` instead of `HTML`, when `Query` was not empty.

## 0.33.0
 - added hash value, which represents the version of the document in `Result.Meta.DocVersion`.

## 0.32.0
 - added support for Russian tags.

## 0.31.0
 - better handling of the page titles for HTML content types;
 - more informative output in verbose mode in CLI - added "title" and "content-type".

## 0.30.0
 - BREAKING CHANGE: changed shape of returned struct by `tagify#Run`;
 - optimized the page title selection.

## 0.29.0
 - BREAKING CHANGE: `tagify#Run` now returns a struct with an extra data along with tags.

## 0.28.0
 - better handling for the apostrophe stopwords.

## 0.27.0
 - `#sanitize` now splits word into parts if there are an unallowed symbols in it;
 - moved onto table-driven tests for same cases.

## 0.26.0
 - simplified sentence splitting regex to only split by either of `.,!;:` symbols.

## 0.25.0
 - improved HTML parser, now it keeps crawling inside the tag even if there are other tags inside;
 - added `<code>` tag for HTML processor;
 - TF part of the TF-IDF score caluclation is now logarithmically scaled;
 - added more tests, plus a bit of a code re-shuffling;
 - bumped `github.com/zoomio/stopwords` to `0.3.0`.

## 0.24.0
 - use TF-IDF for better tags scoring;
 - bumped `github.com/zoomio/inout` to `0.8.2`;
 - improved stop words detection by checking after sanitisation;
 - added `<li>` tag for HTML processor.

## 0.23.0
 - HTML: do not count page's title twice, if it is represented in one of the headings.

## 0.22.1
 - fixed test.

## 0.22.0
 - added `<title>` HTML tag handling;
 - minor code optimizations.

## 0.21.0
 - bumped `github.com/zoomio/inout` to `0.8.0`;
 - removed dependency on `hizer` so now HTML parser actually works better.

## 0.20.0
 - bumped `github.com/zoomio/inout` to `0.7.0`, to handle bigger lines of text;
 - made private internal housekkeping struct `in`.

## 0.19.0
 - bumped Go to `1.13.x`;
 - better error printing/wrapping with `fmt.Fprintf(os.Stderr, ...)` and `fmt.Errorf("...%W...", ...)`;
 - added benchmark test and profiling;
 - improved infrastructure scripts for build and install, added `Makefile`;
 - better help on `--help` option.

## 0.18.0
 - bumped `github.com/zoomio/inout` to `0.6.0`;
 - re-used "self-referential functions and the design of options" approach by Rob Pike by introducing `Option` and new API method `#Run`.

## 0.17.0
 - CSS query (`-q` option): improved overall querying logic, now it retrieves all texts from the matching tags;
 - bumped `github.com/zoomio/inout` to `0.5.0`.

## 0.16.0
 - bumped `github.com/zoomio/inout` to be able to wait for DOM elements to be visible on the web-pages;
 - introduced `-q` option, which is short for "query" to allow to provide CSS query.

## 0.15.0
 - added support for more HTML tags `<h5>`, `<h6>` and `<a>`.

## 0.14.0
 - breaking change: signature of `processor#ParseHTML` changed, removed `bool` argument - `doTagify`, previously it returned a tuple `([]string, []*Tag)` and now it is a single result - `[]*Tag`;
 - increased code coverage.

## 0.13.0
 - breaking change: signature of `processor#ParseHTML` changed, added extra `bool` argument - `tagify` (if set to true, then output `[]*Tag` slice will be populated, otherwise it will be empty) and return tuple values swapped places - `([]string, []*Tag)` instead of `([]*Tag, []string)`.

## 0.12.0
 - bumped `github.com/zoomio/inout` from `0.1.0` to `0.2.0`.

## 0.11.0
 - externalized inout into standalone package `github.com/zoomio/inout`.

## 0.10.0
 - skip words that start with hyphen;
 - moved to Go Modules.

## 0.9.0
 - externalized stop-words into standalone package `github.com/zoomio/stopwords`.

## 0.8.0
 - more refactorring;
 - enabled skipped test;
 - added `-no-stop` boolean flag, to allow disabling of stop-words filter.

## 0.7.0
 - moved stop words in `*.go` file;
 - removed dependency on `github.com/gobuffalo/packr`.

## 0.6.0
 - improved stop words list;
 - changed math for HTML tag weights;
 - improved `#normalize` to be mored defensive in case if sanitize regex still returns not a word.

## 0.5.0
 - removed default and max limits for the tags query;
 - moved `_scripts` to `_bin`;
 - moved `_files` to `_resources`.

## 0.4.0
 - refactored everything; 
 - added comments in some places to better understand logic;
 - added tests for normalization/de-duping;
 - added `-d` option to return tags along with detailed information.

## 0.3.0
 - added de-duplication algorithm based on the inflection;
 - `#GetTags` and `#GetTagsFromString` now accept `contentType` of type `ContentType`, which is more typo-proof.

## 0.2.0
 - code refactoring;
 - simplified internal structure for ease of API use.

## 0.1.0
 - better error handling;
 - code refactoring;
 - `#GetTags` now accepts `int` variable `conteType`;
 - added `#GetTagsFromString` with equal signature to `#GetTags`.

## 0.0.1
 - first release.