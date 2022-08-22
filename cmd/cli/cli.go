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

	source = flag.String("s", "", "source, could be URL (e.g. http://... and https://...) or file path")
	lang   = flag.String("lang", "", "language of the source, e.g. \"en\"")

	// headless
	query = flag.String("q", "", "DOM CSS query, e.g. `-q p` will fetch contents of all <p> tags from the given source")
	ready = flag.String("r", "", "DOM CSS query, waits until certain element available, but fetches contents of the whole HTML document")
	until = flag.Duration("u", 0, "duration to wait before getting HTML contents, handy for SPAs, because they keep loading in browsers for some time")
	img   = flag.String("i", "", "enables capturing screenshot in the provided path")

	limit       = flag.Int("l", 5, "number of tags to return")
	verbose     = flag.Bool("v", false, "enables verbose mode")
	contentType = flag.String("t", tagify.Unknown.String(), fmt.Sprintf("content type of the source, allowed values: %s", strings.Join(config.ContentTypes[:], ", ")))
	noStopWords = flag.Bool("no-stop", true, "removes stop-words from results (see https://github.com/zoomio/stopwords)")
	contentOnly = flag.Bool("content", true, "tagify only content")

	// weighing
	tagWeights          = flag.String("tag-weights", "", "string with the custom tag weights for HTML & Markdown tagging in the form of <tag1>:<score1>|<tag2>:<score2>")
	tagWeightsJSON      = flag.String("tag-weights-json", "", "JSON file with the custom tag weights for HTML & Markdown tagging in the form of { \"<tag1>\": <score1>, \"<tag2>\": <score2> }")
	adjustScores        = flag.Bool("adjust-scores", false, "adjusts tags score to the interval 0.0 to 1.0")
	extraTagWeights     = flag.String("extra-tag-weights", "", "string with the additional tag weights for HTML & Markdown tagging in the form of <tag1>:<score1>|<tag2>:<score2>")
	extraTagWeightsJSON = flag.String("extra-tag-weights-json", "", "JSON file with the additional tag weights for HTML & Markdown tagging in the form of { \"<tag1>\": <score1>, \"<tag2>\": <score2> }")

	// EXPERIMENTAL
	fullSite = flag.Bool("site", false, "[EXPERIMENTAL] might not be included in next releases: allows to tagify full site (HTML only)")

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

	options := []tagify.Option{
		tagify.TargetType(tagify.ContentTypeOf(*contentType)),
		tagify.Limit(*limit),
	}
	if *source != "" {
		options = append(options, tagify.Source(*source))
	}
	if *lang != "" {
		options = append(options, tagify.Language(*lang))
	}

	// headless
	if len(*query) > 0 {
		options = append(options, tagify.Query(*query))
	}
	if len(*ready) > 0 {
		options = append(options, tagify.WaitFor(*ready))
	}
	if *until > 0 {
		options = append(options, tagify.WaitUntil(*until))
	}
	if len(*img) > 0 {
		options = append(options, tagify.Screenshot(true))
	}

	if *verbose {
		options = append(options, tagify.Verbose(*verbose))
	}
	if *noStopWords {
		options = append(options, tagify.NoStopWords(*noStopWords))
	}
	if *contentOnly {
		options = append(options, tagify.ContentOnly(*contentOnly))
	}
	if *fullSite {
		options = append(options, tagify.FullSite(*fullSite))
	}
	if *tagWeights != "" {
		options = append(options, tagify.TagWeightsString(*tagWeights))
	} else if *tagWeightsJSON != "" {
		options = append(options, tagify.TagWeightsJSON(*tagWeightsJSON))
	}
	if *adjustScores {
		options = append(options, tagify.AdjustScores(*adjustScores))
	}
	if *extraTagWeights != "" {
		options = append(options, tagify.ExtraTagWeightsString(*extraTagWeights))
	} else if *extraTagWeightsJSON != "" {
		options = append(options, tagify.ExtraTagWeightsJSON(*extraTagWeightsJSON))
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

	if len(*img) > 0 && len(res.Meta.Screenshot) > 0 {
		err = os.WriteFile(*img, res.Meta.Screenshot, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to store captured screenshot at %s: %v\n", *img, err)
			os.Exit(3)
		}
	}

	if res.RawLen() == 0 {
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
