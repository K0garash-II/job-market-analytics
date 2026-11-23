package extractor

import (
	"strings"
	"unicode"
)

type SkillExtractor struct{ index map[string]int }

func New(index map[string]int) *SkillExtractor { return &SkillExtractor{index} }
func (e *SkillExtractor) Extract(title, description string) []int {
	words := strings.FieldsFunc(strings.ToLower(title+" "+description), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '+' && r != '#' })
	seen := map[int]bool{}
	var ids []int
	for _, word := range words {
		if id, ok := e.index[word]; ok && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}
