package interpol

import (
	"strconv"
	"testing"
)

func TestCounter(t *testing.T) {
	helperTestSingle(t, []testdataSingle{
		{"counter normal", "{{counter, min=0, max=7, step=2}}",
			[]string{"0", "2", "4", "6"}},
		{"counter format test", "{{counter, min=2, max=8, step=2, format=%03d}}",
			[]string{"002", "004", "006", "008"}},
		{"counter backwardtest", "{{counter, min=5, max=0, step=-1}}",
			[]string{"5", "4", "3", "2", "1", "0"}},
		{"counter both", "{{counter, min=1, max=2, step=1}}{{counter, min=7, max=6, step=-1}}",
			[]string{"17", "27", "16", "26"}},
	})
}

func TestRandomNormal(t *testing.T) {
	rnds, err := helperInterpolateAll("{{random, min=0, max=100, count=1000}}")
	if err != nil {
		t.Fatalf("failed to add normal string: %v", err)
	}

	for _, s := range rnds[0] {
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 || n > 100 {
			t.Errorf("random err: min=%d max=%d err=%v random=%s", 0, 100, err, s)
		}
	}
}
