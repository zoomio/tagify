package main

import (
	"context"
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
	doFiltering := flag.Bool("no-stop", true, "filters out all stop-words from results (see https://github.com/zoomio/stopwords)")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(3)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in profiling: %v\n", err)
			os.Exit(3)
		}
		defer pprof.StopCPUProfile()
	}

	options := []tagify.Option{
		tagify.TargetType(tagify.ContentTypeOf(*contentType)),
		tagify.Limit(*limit),
	}
	if *source != "" {
		options = append(options, tagify.Source(*source))
	}
	if *query != "" {
		options = append(options, tagify.Query(*query))
	}
	if *verbose {
		options = append(options, tagify.Verbose(*verbose))
	}
	if *doFiltering {
		options = append(options, tagify.NoStopWords(*doFiltering))
	}

	res, err := tagify.Run(context.Background(), options...)
	if err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "failed to get tags: %v\n", err)
		}
		os.Exit(2)
	}

	if res.Len() == 0 {
		fmt.Println("found 0 tags")
		os.Exit(0)
	}

	if *verbose {
		fmt.Printf("title: %s\n", res.Meta.DocTitle)
		fmt.Printf("content-type: %s\n", res.Meta.ContentType)
	}

	fmt.Printf("%s\n", strings.Join(res.TagsStrings(), " "))
}
