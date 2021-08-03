package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/zoomio/tagify"
	"github.com/zoomio/tagify/config"
)

var (
	version = "tip"

	source         = flag.String("s", "", "source, could be URL (e.g. http://... and https://...) or file path")
	query          = flag.String("q", "", "DOM CSS query, e.g. `-q p` will fetch contents of all <p> tags from the given source")
	limit          = flag.Int("l", 5, "number of tags to return")
	verbose        = flag.Bool("v", false, "enables verbose mode")
	contentType    = flag.String("t", config.Unknown.String(), "type of content type in the source (Text or HTML)")
	noStopWords    = flag.Bool("no-stop", true, "removes stop-words from results (see https://github.com/zoomio/stopwords)")
	tagWeights     = flag.String("tag-weights", "", "string with the custom tag weights for HTML & Markdown tagging in the form of <tag1>:<score1>|<tag2>:<score2>")
	tagWeightsJSON = flag.String("tag-weights-json", "", "JSON file with the custom tag weights for HTML & Markdown tagging in the form of { \"<tag1>\": <score1>, \"<tag2>\": <score2> }")
	adjustScores   = flag.Bool("adjust-scores", false, "adjusts tags score to the interval 0.0 to 1.0")

	// EXPERIMENTAL
	contentOnly = flag.Bool("content", false, "[EXPERIMENTAL] might not be included in next releases: ignore all none content related parts of the page (HTML only)")
	fullSite    = flag.Bool("site", false, "[EXPERIMENTAL] might not be included in next releases: allows to tagify full site (HTML only)")

	// Utility
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	ver        = flag.Bool("version", false, "prints version of Tagify")
)

func main() {
	flag.Parse()

	if *ver {
		fmt.Println(version)
		return
	}

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

	options := []config.Option{
		config.TargetType(config.ContentTypeOf(*contentType)),
		config.Limit(*limit),
	}
	if *source != "" {
		options = append(options, config.Source(*source))
	}
	if *query != "" {
		options = append(options, config.Query(*query))
	}
	if *verbose {
		options = append(options, config.Verbose(*verbose))
	}
	if *noStopWords {
		options = append(options, config.NoStopWords(*noStopWords))
	}
	if *contentOnly {
		options = append(options, config.ContentOnly(*contentOnly))
	}
	if *fullSite {
		options = append(options, config.FullSite(*fullSite))
	}
	if *tagWeights != "" {
		options = append(options, config.TagWeightsString(*tagWeights))
	} else if *tagWeightsJSON != "" {
		options = append(options, config.TagWeightsJSON(*tagWeightsJSON))
	}
	if *adjustScores {
		options = append(options, config.AdjustScores(*adjustScores))
	}

	// print progress updates to terminal
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	if !*verbose {
		wg.Add(1)
		go printProgress(stopCh, &wg)
	}

	res, err := tagify.Run(context.Background(), options...)
	close(stopCh)
	wg.Wait()
	if err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "failed to get tags: %v\n", err)
		}
		os.Exit(2)
	}

	if res.Len() == 0 {
		fmt.Println("found 0 tags")
		return
	}

	if *verbose {
		fmt.Printf("title: %s\n", res.Meta.DocTitle)
		fmt.Printf("hash: %s\n", res.Meta.DocHash)
		fmt.Printf("content-type: %s\n", res.Meta.ContentType)
		println()
	}

	prfx := ""
	if !*verbose {
		prfx = "\r"
	}

	fmt.Fprintf(os.Stdout, "%s%s\n", prfx, strings.Join(res.TagsStrings(), " "))
}

func printProgress(stopCh chan struct{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(80 * time.Millisecond)
	i := -1
	symbs := []string{"|", "\\", "-", "/"}
	for {
		select {
		case <-stopCh:
			wg.Done()
			return
		case <-ticker.C:
			i++
			if i >= len(symbs) {
				i = 0
			}
			fmt.Fprintf(os.Stdout, "\rprocessing... %s", symbs[i])
		}
	}
}
