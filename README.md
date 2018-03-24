# tagify

Gets STDIN, file or HTTP address as an input and returns ordered list of most frequent tags as an output.

Example, get 10 most frequent words from StackOverflow main page:
```bash
$ tagify -s=https://stackoverflow.com -l=10
application using page add file server run ionic local error
```

## Development

* [Go](https://golang.org/dl/)
* [Dep](https://golang.github.io/dep/docs/installation.html)

## TODO

* Introduce verbose mode with self-explanatory output
* Add automatic release pipeline
* Increase test coverage