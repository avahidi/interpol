package test

import (
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
		{"counter both", "{{counter, min=1, max=2, step=1}}{{counter, min=3, max=2, step=-1}}",
			[]string{"13", "23", "12", "22"}},
	})
}
