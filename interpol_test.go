package interpol

import (
	"testing"
)

// testdata
var simpleSingleTestdata = []testdataSingle{
	{"basic  static", "thisthat", []string{"thisthat"}},
	{"basic mixed", "prefix{{counter, min=0, max=10, step=3}}postfix",
		[]string{"prefix0postfix", "prefix3postfix", "prefix6postfix", "prefix9postfix"}},
	{"basic multi",
		"{{counter, min=0, max=2}}{{counter, min=4, max=5}}{{counter, min=9, max=9}}",
		[]string{"049", "149", "249", "059", "159", "259"}},

	{"basic copy",
		"{{counter min=0, max=2 name=counter1}}-{{copy from=counter1}}",
		[]string{"0-0", "1-1", "2-2"}},
	{"double copy",
		"{{counter min=0, max=2 name=counter1}}{{counter min=5, max=6 name=counter2}}-{{copy from=counter1}}{{copy from=counter1}}{{copy from=counter2}}",
		[]string{"05-005", "15-115", "25-225", "06-006", "16-116", "26-226"}},
	{"modifier 1", "{{set data=aBc modifier=toupper}}", []string{"A", "B", "C"}},
	{"modifier 2", "{{set data=aB modifier=toupper}}{{set data=XY modifier=tolower}}",
		[]string{"Ax", "Bx", "Ay", "By"}},
}

var simpleDualTestdata = []testdataDual{
	{"two static", "this", "that", []string{"this"}, []string{"that"}},
	{"two mixed", "this", "{{counter min=0 max=1}}", []string{"this", "this"}, []string{"0", "1"}},
	{"two both", "{{counter min=3 max=2 step=-1}}", "{{counter min=0 max=1}}",
		[]string{"3", "2", "3", "2"}, []string{"0", "0", "1", "1"}},
	{"modifier dual",
		"{{set data=aB modifier=toupper}}", "{{set data=W1 modifier=tolower}}",
		[]string{"A", "B", "A", "B"}, []string{"w", "w", "1", "1"}},
}

func NilCreator(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	return nil, nil
}

// test corner cases and partial use
func TestAddHandler(t *testing.T) {
	ip := New()
	err := ip.AddHandler("monkey", NilCreator)
	if err != nil {
		t.Errorf("could nto add dummy handler, %v", err)
	}

	err = ip.AddHandler("monkey", NilCreator)
	if err == nil {
		t.Errorf("should not add handler twice")
	}
}

func TestInterpolEmpty(t *testing.T) {
	ip := New()
	if _, err := ip.Add(""); err == nil {
		t.Errorf("should not add empty handler")
	}
}

func TestInterpolNone(t *testing.T) {
	ip := New()
	if ip.Next() {
		t.Errorf("terminate when no interpolations exist")
	}
}

func TestInterpolReset(t *testing.T) {
	ip := New()
	str, err := ip.Add("{{counter, min=0, max=4, step=1}}")
	if err != nil {
		t.Errorf("failed to add normal string: %v", err)
	}

	ip.Next() // 1
	ip.Next() // 2
	if str.String() != "2" {
		t.Errorf("incorrect value before reset: %v", str)
	}

	ip.Reset() // 0
	ip.Next()  // 1
	if str.String() != "1" {
		t.Errorf("incorrect value after reset: %v", str)
	}
}

// test normal use cases
func TestInterpolSingle(t *testing.T) {
	helperTestSingle(t, simpleSingleTestdata)
}

func TestInterpolDual(t *testing.T) {
	helperTestDual(t, simpleDualTestdata)
}
