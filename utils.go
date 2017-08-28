package interpol

import (
	"testing"
)

func countOccurrence(list []string, str string) int {
	n := 0
	for _, s := range list {
		if s == str {
			n++
		}
	}
	return n
}

func testGetAll(ip *Interpol, c *InterpolatedString) []string {
	ret := make([]string, 0)
	for {
		ret = append(ret, c.String())
		if !ip.Next() {
			break
		}
	}
	return ret
}

func testCompareAll(t *testing.T, message string, wanted []string, got []string) {
	if len(wanted) != len(got) {
		t.Errorf("%s: wanted %d elements, got %d", message, len(wanted), len(got))
		return
	}

	for i, _ := range got {
		if wanted[i] != got[i] {
			t.Errorf("%s: element %d: wanted '%s' got '%s'", message, i, wanted[i], got[i])
		}
	}
}

func testOccurrence(t *testing.T, message string, valid []string,
	min int, max int, count int, got []string) {
	if len(got) != count {
		t.Errorf("%s: wanted %d elements, got %d", message, count, len(got))
		return
	}

	for _, v := range got {
		n := countOccurrence(valid, v)
		if n < min || n > max {
			t.Errorf("%d occurrences of %s, expected %d-%d", n, v, min, max)
		}
	}
}
