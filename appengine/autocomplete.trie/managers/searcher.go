package managers

import (
	"regexp"
	"strings"
)

type Result struct {
	Query    string    `json:"query"`
	Keywords []Keyword `json:"keywords"`
}

type Keyword struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type Searcher interface {
	Search(string, int) (*Result, error)
}

var productPattern = regexp.MustCompile(`[^\p{L}\p{N}]+`)

func FormatProductKey(prefix, replacement string) string {
	return productPattern.ReplaceAllString(strings.ToLower(strings.TrimSpace(prefix)), replacement)
}
