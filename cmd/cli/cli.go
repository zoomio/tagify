package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/zoomio/tagify"
)

func main() {
	source := flag.String("s", "", "Source")
	limit := flag.Int("l", 0, "Tags limit")
	flag.Parse()

	tagify.InitStopWords()
	fmt.Printf("%v\n", strings.Join(tagify.Process(*source, *limit), " "))
}
