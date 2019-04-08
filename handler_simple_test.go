package interpol

import (
	"testing"
)

func TestCopy(t *testing.T) {
	helperTestSingle(t, []testdataSingle{
		{"copy counter", "{{counter, min=0, max=2 name=cnt}}-{{copy from=cnt}}",
			[]string{"0-0", "1-1", "2-2"}},
	})
}

func TestCopyCorner(t *testing.T) {
	ip := New()
	_, err := ip.Add("{{copy from=doesnexist}}")
	if err == nil {
		t.Errorf("copy target doesnt exist")
	}
}
func TestText(t *testing.T) {
	helperTestSingle(t, []testdataSingle{
		{"text", "sometext",
			[]string{"sometext"}},
		{"text+counter", "before{{counter, min=0, max=2 name=cnt}}after",
			[]string{"before0after", "before1after", "before2after"}},
	})

	h1, err := newTextHandler(nil, "text", nil)
	if err != nil {
		t.Errorf("text handler internal error")
	}
	if h1.String() != "text" {
		t.Errorf("text handler data error")
	}
	if h1.Next() {
		t.Errorf("text handler next error")
	}
}
