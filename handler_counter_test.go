package interpol

import (
	"testing"
)

func TestCounterNormal(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("{{counter, min=0, max=7, step=2}}")
	if err != nil {
		t.Fatalf("failed to add normal string: %v", err)
	}

	testCompareAll(t, "counter norml test",
		[]string{"0", "2", "4", "6"},
		testGetAll(ip, c))
}

func TestCounterFormat(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("{{counter, min=2, max=8, step=2, format=%03d}}")
	if err != nil {
		t.Fatalf("failed to add normal string: %v", err)
	}

	testCompareAll(t, "counter norml test",
		[]string{"002", "004", "006", "008"},
		testGetAll(ip, c))
}

func TestCounterBackward(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("{{counter, min=5, max=0, step=-1}}")
	if err != nil {
		t.Fatalf("failed to add coun tbackward string: %v", err)
	}

	testCompareAll(t, "counter backwardtest",
		[]string{"5", "4", "3", "2", "1", "0"},
		testGetAll(ip, c))
}
