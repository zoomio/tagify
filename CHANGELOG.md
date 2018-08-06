# Changelog
All notable changes to this project will be documented in this file.

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