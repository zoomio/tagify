package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/zoomio/tagify"
)

func main() {
	source := flag.String("s", "", "source, could be URL (e.g. http://... and https://...) or file path")
	query := flag.String("q", "", "DOM CSS query, e.g. `-q p` will fetch contents of all <p> tags from the given source")
	limit := flag.Int("l", 5, "number of tags to return")
	verbose := flag.Bool("v", false, "enables verbose mode")
	contentType := flag.String("t", tagify.Unknown.String(), "type of content type in the source (Text or HTML)")
	doFiltering := flag.Bool("no-stop", true, "enables filtering out stop-words from results")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(3)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	tags, err := tagify.GetTagsWithQuery(*source, *query, tagify.ContentTypeOf(*contentType), *limit, *verbose, *doFiltering)
	if err != nil && *verbose {
		fmt.Fprintf(os.Stderr, "failed to get tags: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("%s\n", strings.Join(tagify.ToStrings(tags), " "))
}
