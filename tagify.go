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

	cfg := config.New(options...)

	var in in
	var err error

	if cfg.Content != "" {
		in = newInFromString(cfg.Content, cfg.ContentType)
	} else {
		in, err = newIn(ctx, cfg)
		if cfg.ContentType > Unknown {
			in.ContentType = cfg.ContentType
		}
	}

	if err != nil {
		return nil, err
	}

	res := processInput(&in, cfg)

	if len(res.RawTags) > 0 {
		if cfg.Verbose {
			fmt.Println("tagifying...")
		}
		res.Tags = processor.Run(cfg, res.Flatten())
		if cfg.Verbose {
			fmt.Printf("\n%v\n", res.Tags)
		}
	}

	return res, nil
}

func processInput(in *in, c *Config) *model.Result {
	switch in.ContentType {
	case HTML:
		res := html.ProcessHTML(c, in)
		if c.Screenshot && len(in.reader.ImgBytes) > 0 {
			res.Meta.Screenshot = in.reader.ImgBytes
		}
		return res
	case Markdown:
		return md.ProcessMD(c, in)
	default:
		return text.ProcessText(c, in)
	}
}
