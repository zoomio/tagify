package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zoomio/tagify"
)

func main() {
	source := flag.String("s", "", "Source")
	query := flag.String("q", "", "DOM CSS query, e.g. `-q p` will fetch contents of all <p> tags from the given source")
	limit := flag.Int("l", 5, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	contentType := flag.String("t", tagify.Unknown.String(), "Content type (Text or HTML)")
	doFiltering := flag.Bool("no-stop", true, "Filter by stop-words")
	flag.Parse()

	tags, err := tagify.GetTagsWithQuery(*source, *query, tagify.ContentTypeOf(*contentType), *limit, *verbose, *doFiltering)
	if err != nil && *verbose {
		println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", strings.Join(tagify.ToStrings(tags), " "))
}
