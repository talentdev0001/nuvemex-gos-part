package part

import "github.com/Montrealist-cPunto/goseanto"

type SearchConfig struct {
	CacheDuration string
	CrawlPool     string
	URL           string
}

type HinterOptions = goseanto.HinterOptions
type HintResult = goseanto.HintResult
