package interpol

import (
	"testing"
)

// test data
var file1 = []string{"first", "second", "third"}      // youtube comment section
var file2 = []string{"Alphonse", "Gabriel", "Capone"} // good old Al

var set1 = []string{"AB", "C", "DE", "XYZ"}
var set2 = []string{"1", "9", "6", "7"}

type randomTestdata struct {
	log             string
	cmd             string
	min, max, count int
	valid           []string
}

var linearFileTestdata = []testdataSingle{
	{"file all-linear", "{{file, filename=testdata/data1.txt mode=linear count=0}}", file1},
	{"file some-linear", "{{file, filename=testdata/data1.txt mode=linear count=2}}", file1[:2]},
	{"file more-linear", "{{file, filename=testdata/data1.txt mode=linear count=5}}",
		append(file1, file1[0], file1[1])},
	{"file all-linear-empty", "{{file, filename=testdata/data1.txt mode=linear empty=true}}", append(file1, "")},
	{"file more-linear-empty", "{{file, filename=testdata/data1.txt mode=linear empty=true count=5}}",
		append(file1, "", file1[0])},
}

var randomFileTests = []randomTestdata{
	{"file perm-all", "{{file, filename=testdata/data1.txt mode=perm count=0}}", 1, 1, 3, file1},
	{"file perm-some", "{{file, filename=testdata/data1.txt mode=perm count=2}}", 1, 1, 2, file1},
	{"file perm-many", "{{file, filename=testdata/data1.txt mode=perm count=5}}", 1, 2, 5, file1},

	{"file rand-all", "{{file, filename=testdata/data2.txt mode=rand count=0}}", 0, 3, 3, file2},
	{"file rand-all", "{{file, filename=testdata/data2.txt mode=rand count=2}}", 0, 2, 2, file2},
	{"file rand-all", "{{file, filename=testdata/data2.txt mode=rand count=6}}", 0, 6, 6, file2},
}

var linearSetTestdata = []testdataSingle{
	// linear tests
	{"set linear-all-chars", "{{set data=1967 mode=linear count=0}}", set2},
	{"set linear-all", "{{set data=AB;C;DE;XYZ sep=; mode=linear count=0}}", set1},
	{"set linear-some", "{{set data=AB;C;DE;XYZ sep=; mode=linear count=2}}", set1[:2]},
	{"set linear-many", "{{set data=AB;C;DE;XYZ sep=; mode=linear count=5}}", append(set1, set1[0])},
}

var randomSetTests = []randomTestdata{
	{"set perm-all", "{{set data=1967 mode=perm count=0}}", 1, 1, 4, set2},
	{"set perm-some", "{{set data=1967 mode=perm count=2}}", 1, 1, 2, set2},
	{"set perm-many", "{{set data=1967 mode=perm count=5}}", 1, 2, 5, set2},
	{"set rand-all", "{{set data=1967 mode=rand count=0}}", 0, 3, 4, set2},
	{"set rand-all", "{{set data=1967 mode=rand count=2}}", 0, 2, 2, set2},
	{"set rand-all", "{{set data=1967 mode=rand count=6}}", 0, 6, 6, set2},
}

// helpers

func testPermutationRepeat(t *testing.T, log string, cmd string, repeats int, setSize int) {
	for n := 0; n < repeats; n++ {
		strs, err := helperInterpolateAll(cmd)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", log, err)
		}

		got := strs[0][setSize-1 : setSize+1]
		if got[0] == got[1] {
			t.Errorf("%s: permutation repeats during wrap-around: %v", log, got)
			return
		}
	}
}

// file tests
func TestFileLinear(t *testing.T) {
	helperTestSingle(t, linearFileTestdata)
}

func TestFileRandom(t *testing.T) {
	for _, test := range randomFileTests {
		helperTestRange(t, test.log, test.cmd, test.min, test.max, test.count,
			test.valid)
	}
}

func TestFilePermRepeat(t *testing.T) {
	for n := 0; n < 100; n++ {
		strs, err := helperInterpolateAll("{{file, filename=testdata/data1.txt mode=perm count=20}}")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		got := strs[0][len(file1)-1 : len(file2)+1]
		if got[0] == got[1] {
			t.Errorf("permutation repeats two: %v", got)
			return
		}
	}
}

// permutation is not linear
func TestFilePerm(t *testing.T) {
	testPermutationRepeat(t, "file perm wrap", "{{file, filename=testdata/data1.txt mode=perm count=20}}", 100, len(file1))
}

// set
func TestSetinear(t *testing.T) {
	helperTestSingle(t, linearSetTestdata)
}

func TestSetRandom(t *testing.T) {
	for _, test := range randomSetTests {
		helperTestRange(t, test.log, test.cmd, test.min, test.max, test.count,
			test.valid)
	}
}

func TestSetPermRepeat(t *testing.T) {
	testPermutationRepeat(t, "set perm wrap", "{{set data=ABCDE count=20}}", 100, 5)
}
