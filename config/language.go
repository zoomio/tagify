package config

import (
	"fmt"

	"github.com/abadojack/whatlanggo"
	"github.com/zoomio/stopwords"
)

type Vocabulary interface {
	SetLang(l string)
	SetReg(r *stopwords.Register)
}

// DetectLang detects language and setups the stop words for it.
func DetectLang(cfg *Config, controlStr string, contents Vocabulary) {
	if len(cfg.Lang) == 0 {
		info := whatlanggo.Detect(controlStr)
		if cfg.Verbose {
			fmt.Printf("detected language based on %q: %s [%s] [%s], confidence %2.f\n",
				controlStr, info.Lang.String(), info.Lang.Iso6391(), info.Lang.Iso6393(), info.Confidence)
		}
		if info.IsReliable() {
			contents.SetLang(info.Lang.String())
			cfg.SetStopWords(info.Lang.Iso6391())
		} else {
			contents.SetLang("English")
			cfg.SetStopWords("en")
			if cfg.Verbose {
				fmt.Println("use English language hence detection is not reliable")
			}
		}
	} else {
		contents.SetLang(cfg.Lang)
		cfg.SetStopWords(cfg.Lang)
	}
	if cfg.NoStopWords {
		contents.SetReg(cfg.StopWords)
	}

}
