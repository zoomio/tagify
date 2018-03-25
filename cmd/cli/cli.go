package main

import (
	"os"
	"flag"
	"strings"

	"github.com/gobuffalo/packr"

	"github.com/zoomio/tagify"
	"github.com/zoomio/tagify/rank"
)

func main() {
	source := flag.String("s", "", "Source")
	limit := flag.Int("l", 0, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	contentType := flag.Int("t", -1, "Content type")
	flag.Parse()

	box := packr.NewBox("../../_files")

	var err error
	var tags []string
	
	err = rank.InitStopWords(&box)
	if err != nil && *verbose {
		println(err)
		os.Exit(1)
	}

	tags, err = tagify.GetTags(*source, *contentType, *limit, *verbose)
	if err != nil && *verbose {
		println(err)
		os.Exit(2)
	}

	println(strings.Join(tags, " "))
}
