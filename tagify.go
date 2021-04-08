package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/processor"
	"github.com/zoomio/tagify/processor/html"
	"github.com/zoomio/tagify/processor/md"
	"github.com/zoomio/tagify/processor/model"
	"github.com/zoomio/tagify/processor/text"
)

type config struct {
	source      string
	query       string
	content     string
	contentType ContentType
	limit       int
	verbose     bool
	noStopWords bool
	contentOnly bool
	fullSite    bool
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
func ToStrings(items []*model.Tag) []string {
	return model.ToStrings(items)
}

func processInput(in *in, c config) (tags []*model.Tag, pageTitle string, hash []byte) {
	var out *model.ParseOutput
	switch in.ContentType {
	case HTML:
		out = html.ParseHTML(in,
			model.Verbose(c.verbose),
			model.NoStopWords(c.noStopWords),
			model.ContentOnly(c.contentOnly),
			model.FullSite(c.fullSite),
			model.Source(in.source))
	case Markdown:
		out = md.ParseMD(in,
			model.Verbose(c.verbose),
			model.NoStopWords(c.noStopWords))
	default:
		out = text.ParseText(in,
			model.Verbose(c.verbose),
			model.NoStopWords(c.noStopWords))
	}

	pageTitle = out.DocTitle
	hash = out.DocHash

	if len(out.Tags) > 0 {
		if c.verbose {
			fmt.Println("tagifying...")
		}
		tags = processor.Run(out.FlatTags(), c.limit)
		if c.verbose {
			fmt.Printf("\n%v\n", tags)
		}
	}

	return
}
