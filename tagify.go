package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/model"
	"github.com/zoomio/tagify/processor"
	"github.com/zoomio/tagify/processor/html"
	"github.com/zoomio/tagify/processor/md"
	"github.com/zoomio/tagify/processor/text"
)

// Run produces slice of tags ordered by frequency.
func Run(ctx context.Context, options ...Option) (*model.Result, error) {

	c := config.New(options...)

	var in in
	var err error

	if c.Content != "" {
		in = newInFromString(c.Content, c.ContentType)
	} else {
		in, err = newIn(ctx, c.Source, c.Query, c.Verbose)
		if c.ContentType > Unknown {
			in.ContentType = c.ContentType
		}
	}

	if err != nil {
		return nil, err
	}

	res := processInput(&in, c)

	// var tags []*model.Tag
	if len(res.RawTags) > 0 {
		if c.Verbose {
			fmt.Println("tagifying...")
		}
		res.Tags = processor.Run(c, res.FlatTags())
		if c.Verbose {
			fmt.Printf("\n%v\n", res.Tags)
		}
	}

	return res, nil
}

// ToStrings transforms a list of tags into a list of strings.
func ToStrings(items []*model.Tag) []string {
	return model.ToStrings(items)
}

func processInput(in *in, c *Config) *model.Result {
	switch in.ContentType {
	case HTML:
		return html.ParseHTML(c, in)
	case Markdown:
		return md.ParseMD(c, in)
	default:
		return text.ParseText(c, in)
	}
}
