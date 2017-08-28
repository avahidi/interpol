package interpol

import (
	"testing"
)

func NilCreator(text string, data *InterpolatorData) (Handler, error) {
	return nil, nil
}

func TestAddHandler(t *testing.T) {
	ip := NewInterpol()
	err := ip.AddHandler("monkey", NilCreator)
	if err != nil {
		t.Errorf("first handler expected no error, %v", err)
	}

	err = ip.AddHandler("monkey", NilCreator)
	if err == nil {
		t.Errorf("second handler expected error")
	}
}

func TestAddEmpty(t *testing.T) {
	ip := NewInterpol()
	if _, err := ip.Add(""); err == nil {
		t.Errorf("expected error when adding empty string")
	}
}

func TestAddNormal(t *testing.T) {
	ip := NewInterpol()
	str, err := ip.Add("thisthat")
	if err != nil {
		t.Errorf("failed to add normal string: %v", err)
	}

	if str.String() != "thisthat" {
		t.Errorf("incorrect normal string: %v", str)
	}

	if ip.Next() {
		t.Errorf("Expected only one element from static text")
	}
}

func TestAddMixed(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("prefix{{counter, min=0, max=10, step=3}}postfix")
	if err != nil {
		t.Fatalf("failed to add normal string: %v", err)
	}

	testCompareAll(t, "add mixed test",
		[]string{"prefix0postfix", "prefix3postfix", "prefix6postfix", "prefix9postfix"},
		testGetAll(ip, c))
}

func TestAddMulti(t *testing.T) {
	ip := NewInterpol()
	c, err := ip.Add("{{counter, min=0, max=2}}{{counter, min=4, max=5}}{{counter, min=9, max=9}}")
	if err != nil {
		t.Fatalf("failed to add multi string: %v", err)
	}

	testCompareAll(t, "add multi test",
		[]string{"049", "149", "249", "059", "159", "259"},
		testGetAll(ip, c))
}
