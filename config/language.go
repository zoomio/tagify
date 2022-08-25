package config

import (
	"fmt"

	"github.com/abadojack/whatlanggo"
)

// DetectLang detects language and setups the stop words for it.
func DetectLang(cfg *Config, controlStr string) {
	if len(cfg.Lang) == 0 {
		info := whatlanggo.Detect(controlStr)
		if cfg.Verbose {
			fmt.Printf("detected language based on %q: %s [%s] [%s], confidence %2.f\n",
				controlStr, info.Lang.String(), info.Lang.Iso6391(), info.Lang.Iso6393(), info.Confidence)
		}
		if info.IsReliable() {
			setLang(cfg, info.Lang.Iso6391())
		} else {
			setLang(cfg, "en")
		}
	} else {
		setLang(cfg, cfg.Lang)
	}
}

func setLang(cfg *Config, lang string) {
	cfg.Lang = lang
	cfg.SetStopWords(lang)
	if cfg.Verbose {
		fmt.Printf("language to use: %s\n", lang)
	}
}
