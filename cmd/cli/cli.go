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
	limit := flag.Int("l", 5, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	contentType := flag.String("t", tagify.Unknown.String(), "Content type (Text or HTML)")
	doFiltering := flag.Bool("no-stop", true, "Filter by stop-words")
	flag.Parse()

	err := tagify.Init()
	if err != nil && *verbose {
		println(err)
		os.Exit(1)
	}

	cntType := tagify.ContentTypeOf(*contentType)
	tags, err := tagify.GetTags(*source, cntType, *limit, *verbose, *doFiltering)
	if err != nil && *verbose {
		println(err)
		os.Exit(2)
	}

	fmt.Printf("%s\n", strings.Join(tagify.ToStrings(tags), " "))
}
