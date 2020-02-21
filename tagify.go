package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/processor"
)

type config struct {
	source      string
	query       string
	content     string
	contentType ContentType
	limit       int
	verbose     bool
	noStopWords bool
}

// Run produces slice of tags ordered by frequency.
func Run(ctx context.Context, options ...Option) (*Result, error) {

	c := &config{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	var in in
	var err error

	if c.content != "" {
		in = newInFromString(c.content, c.contentType)
	} else {
		in, err = newIn(ctx, c.source, c.query, c.verbose)
		if c.contentType > Unknown {
			in.ContentType = c.contentType
		} else if c.query != "" {
			in.ContentType = HTML
		}
	}

	if err != nil {
		return nil, err
	}

	tags, title, hash := processInput(&in, *c)

	return &Result{
		Meta: &Meta{
			ContentType: in.ContentType,
			DocTitle:    title,
			DocHash:     fmt.Sprintf("%x", hash),
		},
		Tags: tags,
	}, nil
}

// ToStrings transforms a list of tags into a list of strings.
func ToStrings(items []*processor.Tag) []string {
	return processor.ToStrings(items)
}

func processInput(in *in, c config) (tags []*processor.Tag, pageTitle string, hash []byte) {
	var out *processor.ParseOutput
	switch in.ContentType {
	case HTML:
		out = processor.ParseHTML(in, processor.Verbose(c.verbose), processor.NoStopWords(c.noStopWords))
	case Markdown:
		out = processor.ParseMD(in, processor.Verbose(c.verbose), processor.NoStopWords(c.noStopWords))
	default:
		out = processor.ParseText(in, processor.Verbose(c.verbose), processor.NoStopWords(c.noStopWords))
	}

	pageTitle = out.DocTitle
	hash = out.DocHash

	if len(out.Tags) > 0 {
		if c.verbose {
			fmt.Println("tagifying...")
		}
		tags = processor.Run(out.Tags, c.limit)
		if c.verbose {
			fmt.Printf("\n%v\n", tags)
		}
	}

	return
}
