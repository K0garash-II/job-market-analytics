package extractor

import (
	"regexp"
	"strings"
)

type term struct {
	id      int
	pattern *regexp.Regexp
}
type SkillExtractor struct{ terms []term }

func New(index map[string]int) *SkillExtractor {
	e := &SkillExtractor{}
	for name, id := range index {
		if name == "" {
			continue
		}
		pattern := `(^|[^\pL\pN])` + regexp.QuoteMeta(strings.ToLower(name)) + `($|[^\pL\pN])`
		e.terms = append(e.terms, term{id, regexp.MustCompile(pattern)})
	}
	return e
}

func (e *SkillExtractor) Extract(title, description string) []int {
	text := strings.ToLower(title + " " + description)
	var ids []int
	for _, term := range e.terms {
		if term.pattern.MatchString(text) {
			ids = append(ids, term.id)
		}
	}
	return ids
}
