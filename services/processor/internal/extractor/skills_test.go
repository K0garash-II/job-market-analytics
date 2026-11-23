package extractor

import "testing"

func TestExtractWholeWords(t *testing.T) {
	e := New(map[string]int{"go": 1, "java": 2})
	ids := e.Extract("Django, затем Go", "JavaScript GO")
	if len(ids) != 1 || ids[0] != 1 {
		t.Fatal(ids)
	}
}
