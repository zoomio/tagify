package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/processor"
	"github.com/zoomio/tagify/processor/html"
	"github.com/zoomio/tagify/processor/md"
	"github.com/zoomio/tagify/processor/model"
	"github.com/zoomio/tagify/processor/text"
)

// Run produces slice of tags ordered by frequency.
func Run(ctx context.Context, options ...config.Option) (*Result, error) {

	c := config.New(options...)

	var in in
	var err error

	if c.Content != "" {
		in = newInFromString(c.Content, c.ContentType)
	} else {
		in, err = newIn(ctx, c.Source, c.Query, c.Verbose)
		if c.ContentType > config.Unknown {
			in.ContentType = c.ContentType
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

func processInput(in *in, c config.Config) (tags []*model.Tag, pageTitle string, hash []byte) {
	var out *model.ParseOutput

	opts := []model.ParseOption{}
	if c.TagWeights != "" {
		opts = append(opts, model.TagWeightsString(c.TagWeights))
	} else if c.TagWeightsJSON != "" {
		opts = append(opts, model.TagWeightsJSON(c.TagWeightsJSON))
	}

	switch in.ContentType {
	case config.HTML:
		out = html.ParseHTML(&c, in, opts...)
	case config.Markdown:
		out = md.ParseMD(&c, in, opts...)
	default:
		out = text.ParseText(&c, in, opts...)
	}

	pageTitle = out.DocTitle
	hash = out.DocHash

	if len(out.Tags) > 0 {
		if c.Verbose {
			fmt.Println("tagifying...")
		}
		tags = processor.Run(&c, out.FlatTags())
		if c.Verbose {
			fmt.Printf("\n%v\n", tags)
		}
	}

	return
}
