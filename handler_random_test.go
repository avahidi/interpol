package interpol

import (
	"strconv"
	"testing"
)

func TestRandomNormal(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("{{random, min=0, max=100, count=1000}}")
	if err != nil {
		t.Fatalf("failed to add normal string: %v", err)
	}

	rnds := testGetAll(ip, c)
	if len(rnds) != 1000 {
		t.Errorf("Expected %d randoms, got %d\n", 1000, len(rnds))
	}

	for _, s := range rnds {
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 || n > 100 {
			t.Errorf("random err: min=%d max=%d err=%v random=%s", 0, 100, err, s)
		}
	}
}
