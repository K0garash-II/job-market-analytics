package extractor

import (
	"slices"
	"testing"
)

func TestExtractWholeWords(t *testing.T) {
	e := New(map[string]int{"go": 1, "java": 2, "Node.js": 3, "c++": 4})
	ids := e.Extract("Django, затем Go", "JavaScript GO Node.js и C++")
	slices.Sort(ids)
	if !slices.Equal(ids, []int{1, 3, 4}) {
		t.Fatal(ids)
	}
}
