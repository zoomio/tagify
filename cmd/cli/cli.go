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
	limit := flag.Int("l", 0, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	detailed := flag.Bool("d", false, "Detailed")
	contentType := flag.String("t", tagify.Unknown.String(), "Content type (Text or HTML)")
	flag.Parse()

	err := tagify.Init()
	if err != nil && *verbose {
		println(err)
		os.Exit(1)
	}

	tags, err := tagify.GetTags(*source, tagify.ContentTypeOf(*contentType), *limit, *verbose)
	if err != nil && *verbose {
		println(err)
		os.Exit(2)
	}

	if *detailed {
		fmt.Printf("%#v", tags)
		return
	}

	println(strings.Join(tagify.ToStrings(tags), " "))
}
