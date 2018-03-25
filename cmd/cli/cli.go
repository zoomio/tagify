package main

import (
	"os"
	"flag"
	"strings"
	
	"github.com/zoomio/tagify"
)

func main() {
	source := flag.String("s", "", "Source")
	limit := flag.Int("l", 0, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	contentType := flag.String("t", tagify.Unknown.String(), "Content type (Text or HTML)")
	flag.Parse()

	err := tagify.Init()
	if err != nil && *verbose {
		println(err)
		os.Exit(1)
	}

	tags, err := tagify.GetTags(*source, *contentType, *limit, *verbose)
	if err != nil && *verbose {
		println(err)
		os.Exit(2)
	}

	println(strings.Join(tags, " "))
}
