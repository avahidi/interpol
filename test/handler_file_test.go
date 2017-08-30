package test

import (
	"testing"
)

// test data
var data1 = []string{"first", "second", "third"}      // youtube comment section
var data2 = []string{"Alphonse", "Gabriel", "Capone"} // good old Al

var linearTestData = []testdataSingle{
	// linear tests
	{"file all-linear", "{{file, filename=testdata1.txt mode=linear count=0}}", data1},
	{"file some-linear", "{{file, filename=testdata1.txt mode=linear count=2}}", data1[:2]},
	{"file more-linear", "{{file, filename=testdata1.txt mode=linear count=5}}",
		append(data1, data1[0], data1[1])},
}

type randomTestdata struct {
	log             string
	cmd             string
	min, max, count int
	valid           []string
}

var randomTests = []randomTestdata{
	{"file perm-all", "{{file, filename=testdata1.txt mode=perm count=0}}", 1, 1, 3, data1},
	{"file perm-some", "{{file, filename=testdata1.txt mode=perm count=2}}", 1, 1, 2, data1},
	{"file perm-many", "{{file, filename=testdata1.txt mode=perm count=5}}", 1, 2, 5, data1},

	{"file rand-all", "{{file, filename=testdata2.txt mode=rand count=0}}", 0, 3, 3, data2},
	{"file rand-all", "{{file, filename=testdata2.txt mode=rand count=2}}", 0, 2, 2, data2},
	{"file rand-all", "{{file, filename=testdata2.txt mode=rand count=6}}", 0, 6, 6, data2},
}

// test linear
func TestFileLinear(t *testing.T) {
	helperTestSingle(t, linearTestData)
}

// testing permutations and random
func TestFileRandom(t *testing.T) {
	for _, test := range randomTests {
		helperTestRange(t, test.log, test.cmd, test.min, test.max, test.count,
			test.valid)
	}
}

// test permutation wrap-around
func TestFilePermRepeat(t *testing.T) {
	for n := 0; n < 100; n++ {
		strs, err := helperInterpolateAll("{{file, filename=testdata1.txt mode=perm count=20}}")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		got := strs[0][len(data1)-1 : len(data1)+1]
		if got[0] == got[1] {
			t.Errorf("permutation repeats two: %v", got)
			return
		}
	}
}

// permutation is not linear
func TestFilePerm(t *testing.T) {
	for n := 0; n < 100; n++ {
		strs, err := helperInterpolateAll("{{file, filename=testdata1.txt mode=perm}}")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		for i, got := range strs[0] {
			if got != data1[i] {
				return
			}
		}
	}

	t.Errorf("Permutation doesnt change order")
}
