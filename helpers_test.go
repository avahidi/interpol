package interpol

import (
	"fmt"
	"testing"
)

type testdataSingle struct {
	log      string
	cmd      string
	expected []string
}

type testdataDual struct {
	log       string
	cmd1      string
	cmd2      string
	expected1 []string
	expected2 []string
}

// create interpolation and get all data
func helperInterpolateAll(cmds ...string) ([][]string, error) {
	ip := New()
	strs, err := ip.AddMultiple(cmds...)
	if err != nil {
		return nil, err
	}
	return helperExtractAll(ip, strs...), nil
}

// get all data from already created interpolations
func helperExtractAll(ip *Interpol, strs ...*InterpolatedString) [][]string {
	ret := make([][]string, len(strs))
	for ip.Next() {
		for i := range strs {
			ret[i] = append(ret[i], strs[i].String())
		}
	}
	return ret
}

func helperTestSingle(t *testing.T, testdata []testdataSingle) {
	for _, test := range testdata {
		strs, err := helperInterpolateAll(test.cmd)
		if err != nil {
			t.Errorf("%s: %v", test.log, err)
			break
		}

		got := strs[0]
		if len(got) != len(test.expected) {
			t.Errorf("%s: expected %d elements got %d", test.log, len(test.expected), len(got))
			fmt.Printf("%v - %v\n", test.expected, got) // DEBUG
			return
		}
		for i, e := range test.expected {
			if got[i] != e {
				t.Errorf("%s: element %d expected '%s' got '%s'", test.log, i, e, got[i])
			}
		}
	}
}

func helperTestDual(t *testing.T, testdata []testdataDual) {
	for _, test := range testdata {
		strs, err := helperInterpolateAll(test.cmd1, test.cmd2)
		if err != nil {
			t.Errorf("%s: %v", test.log, err)
			break
		}

		got1, got2 := strs[0], strs[1]
		if len(got1) != len(test.expected1) {
			t.Errorf("%s: expected %d elements got %d", test.log, len(test.expected1), len(got1))
			return
		}
		for i := range got1 {
			if got1[i] != test.expected1[i] {
				t.Errorf("%s: first element %d expected '%s' got '%s'",
					test.log, i, test.expected1[i], got1[i])
			}
			if got2[i] != test.expected2[i] {
				t.Errorf("%s: second element %d expected '%s' got '%s'",
					test.log, i, test.expected2[i], got2[i])
			}
		}
	}
}

func helperCountOccurrence(list []string, str string) int {
	n := 0
	for _, s := range list {
		if s == str {
			n++
		}
	}
	return n
}

// file test helper function
func helperTestRange(t *testing.T, log string, cmd string, min int, max int, count int, valid []string) {
	res, err := helperInterpolateAll(cmd)
	if err != nil {
		t.Fatalf("%s: %v", log, err)
	}

	got := res[0]
	if len(got) != count {
		t.Errorf("%s: wanted %d elements, got %d", log, count, len(got))
		return
	}

	for _, v := range got {
		n := helperCountOccurrence(valid, v)
		if n < min || n > max {
			t.Errorf("%d occurrences of %s, expected %d-%d", n, v, min, max)
		}
	}
}
