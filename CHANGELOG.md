# Changelog

## v0.61.0

- bumped Go to 1.20;
- bumped `github.com/zoomio/inout` to v0.13.0.

## v0.60.2

- fixed dictionary loader for segmenter for Chinese & Japanese languages.

## v0.60.1

- BREAKING: from now on `ContentOnly` option is set to `true` by default;
- optimization: moved segmenter inside the config with the lazy initialization so now it happens only once;
- fix: in cases when language detection is reliable it is now using correct value;
- fix: use the same segmenter logic in the plain text processor.

## v0.60.0

- graduated `ContentOnly` option (`-content` option in the CLI mode); 
- BREAKING: from now on `-content` option in the CLI mode is set to `true` by default.

## v0.59.0

- use different segmentation logic based on the `github.com/go-ego/gse` segmenter for Chinese & Japanese languages;
- improved HTML parser logic: optimised the way it collects contents of a document and improved logic for splitting into sentences;
- fallback to the English language for the stop words in cases when language detection is not reliable;
- added `lang` option to the CLI to be able to provide the language of the document;
- bumped `github.com/zoomio/stopwords` to `0.11.0`.

## v0.58.0

- stopped ignoring `<h1>` in cases when they are equal to the `<title>`, as in now they are included.

## v0.57.0

- Bumped `github.com/zoomio/inout` to `0.12.0`;
- Fixed `-q` option or `Query` in the code (HTTP/HTML mode only), so now it actually works and retrieves contents of the DOM element for the query;
- Introduced `-r` option or `WaitFor` (HTTP/HTML mode only) to allow for waiting for certain DOM element to be ready before getting HTML;
- Introduced `-u` option or `WaitUntil` (HTTP/HTML mode only) to allow to wait for a certain delay before getting HTML;
- Introduced `-i` option or `Screenshot` (HTTP/HTML mode only) to capture a full screenshot of HTML in the given path.

## v0.56.1

- Added macOS (darwin) ARM64 release.

## v0.56.0

- Bumped Go to 1.18;
- BREAKING: renamed `ParseHTML`, `ParseMD` & `ParseText` to `ProcessHTML`, `ProcessMD` & `ProcessText` respectively;
- BREAKING: renamed `extension.Result` to `extension.ExtResult`;
- New option `AllTagWeights` for enabling parsing through everything;
- New option `ExcludeTagsString` for prohibitting some of the tags;
- `ParseHTML` & `ParseMD` are made public to open up parsing capabilities.

## v0.55.0

- improved handling of the words with the "`" or "'" symbols.

## v0.54.0

- BREAKING FROM 0.53.0: changed `config.StopWords` option signature to expect a slice of strings instead of `*stopwords.Register`;
- bumped `github.com/zoomio/stopwords` to `0.10.0`.

## v0.53.0

- added `Option` called `StopWords` to allow for custom stop-words setup, also made `Domains` variable public.

## v0.52.0

- added URL sanitization in the texts, so it excludes things like http, https & www from them.

## v0.51.0

- HTML processor: fallback to <h1> tag (if any) in case if the <title> has not been provided for some reason;
- HTML processor: use the longest parsed line in order to detect document language.

## v0.50.0

- [BREAKING CHANGE (most likely)] extensions (BETA) release - this is the BIGGEST RELEASE since the addition of the Markdown (documentation is in progress);

## v0.49.2

- support backwards compatibility for `ContentTypeOf`.

## v0.49.1

- same as `v0.49.0`.

## v0.49.0

- added language detection in order to improve handling of stop words.

## v0.48.0

- FEATURE: added new parameter `-adjust-scores` to allow configuring scores adjustment to the interval from 0.0 to 1.0.

## v0.47.0

- consider only the <title> tags which are part of the <head>.

## v0.46.0

- FEATURE: added two new parameters `-tag-weights` and `-tag-weights-json` to allow configuring parsed tags & weights for HTML and Markdown sources;
- FEATURE: HTML mode is now parsing contents of `<meta name="description" content="...">` by default;
- MISC: re-organised `processor` package into smaller focused sub-packages.

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