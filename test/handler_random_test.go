package test

import (
	"strconv"
	"testing"
)

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
