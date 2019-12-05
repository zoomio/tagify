# Changelog
All notable changes to this project will be documented in this file.

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