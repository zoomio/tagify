package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gobuffalo/packr"

	"github.com/zoomio/tagify"
	"github.com/zoomio/tagify/rank"
)

func main() {
	source := flag.String("s", "", "Source")
	limit := flag.Int("l", 0, "Tags limit")
	verbose := flag.Bool("v", false, "Verbose mode")
	flag.Parse()

	box := packr.NewBox("../../_files")
	rank.InitStopWords(&box)
	fmt.Printf("%v\n", strings.Join(tagify.GetTags(*source, *limit, *verbose), " "))
}
