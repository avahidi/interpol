package interpol

import (
	"testing"
)

// test internal scanner
func TestScanner(t *testing.T) {
	var testdata = []struct {
		input  string
		output []string
	}{
		{"", []string{}},
		{"one", []string{"one"}},
		{"one two", []string{"one", "two"}},
		{"  spacey two", []string{"spacey", "two"}},
		{"with ' some qoutes '", []string{"with", " some qoutes "}},
		{"var=123", []string{"var", "=", "123"}},
		{"var='=123'", []string{"var", "=", "=123"}},
		{"var\\==123", []string{"var\\=", "=", "123"}},
	}

	for _, test := range testdata {
		s := newScanner(test.input)
		for _, str := range test.output {
			str2, typ := s.next()
			if typ == teof {
				t.Errorf("Unexpected EOF in '%s'", test.input)
				break
			}
			if typ == terror {
				t.Errorf("Unexpected Error in '%s'", test.input)
				break
			}
			if str != str2 {
				t.Errorf("Expected '%s' got '%s'", str, str2)
			}
		}

		str, typ := s.next()
		if typ != teof {
			t.Errorf("Expected EOF, got '%s':%d in '%s'", str, typ, test.input)
		}
	}
}

type interpolParserTestdata struct {
	log  string
	cmd  string
	typ  string
	data map[string]string
}

type lineParserTestdata struct {
	log      string
	line     string
	expected []textElement
}

var interpolTestdata = []interpolParserTestdata{
	{"keep space 1", "spacetype datax=' '", "spacetype",
		map[string]string{"datax": " "}},
	{"nodata", "justtype", "justtype", map[string]string{}},

	{"withdata", "type1 data1 data2=value2", "type1",
		map[string]string{"data1": "", "data2": "value2"}},
	{"remove space 1", "spacetype datax  = valuex ", "spacetype",
		map[string]string{"datax": "valuex"}},
	{"remove space 2", "spacetype nodata datax = valuex ", "spacetype",
		map[string]string{"nodata": "", "datax": "valuex"}},
	{"remove space 3", "spacetype nodata datax =valuex ", "spacetype",
		map[string]string{"nodata": "", "datax": "valuex"}},
	{"remove space 4", "spacetype nodata1 nodata2", "spacetype",
		map[string]string{"nodata1": "", "nodata2": ""}},
}

var lineTestdata = []lineParserTestdata{
	{"just text", "justtext", []textElement{
		{static: true, text: "justtext"}}},
	{"just interpol", "{{justinterpol}}", []textElement{
		{static: false, text: "justinterpol"}}},
	{"mixed line", "stuffbefore{{interpol}}stuffafter", []textElement{
		{static: true, text: "stuffbefore"},
		{static: false, text: "interpol"},
		{static: true, text: "stuffafter"}}},
}

// test parsing of interpol commands
func TestParseInterpol(t *testing.T) {
	for _, test := range interpolTestdata {
		id, err := parseInterpolator(test.cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", test.log, err)
		} else {
			if id.Type != test.typ {
				t.Errorf("%s: expected type '%s' got '%s'", test.log, test.typ, id.Type)
			}
			if len(id.Properties) != len(test.data) {
				t.Errorf("%s: expected %d properties, got %d",
					test.log, len(test.data), len(id.Properties))
			} else {
				for k, v := range test.data {
					if id.Properties[k] != v {
						t.Errorf("%s: expected %s='%s', got '%s'",
							test.log, k, v, id.Properties[k])
					}
				}
			}
		}
	}
}

func TestParseInterpolCorner(t *testing.T) {
	_, err := parseInterpolator("")
	if err == nil {
		t.Error("should not allow empty interpol")
	}
}

// test parsing of lines:
func TestParseLine(t *testing.T) {
	for _, test := range lineTestdata {
		els, err := parseLine(test.line)
		if err != nil {
			t.Errorf("%s: could not parse line, %v", test.log, err)
		} else {
			if len(test.expected) != len(els) {
				t.Errorf("%s: expected %d elements, got %d",
					test.log, len(test.expected), len(els))
			} else {
				for i, e := range test.expected {
					if e.static != els[i].static || e.text != els[i].text {
						t.Errorf("%s: element %d expected %v got %v",
							test.log, i, e, els[i])
					}
				}
			}
		}
	}
}

func TestParseLineCorner(t *testing.T) {
	_, err := parseLine("")
	if err == nil {
		t.Error("should not allow empty line")
	}

	_, err = parseLine("{{counter")
	if err == nil {
		t.Error("should not allow open interpolator")
	}
}

func TestRemoveQoute(t *testing.T) {
	var testdata = []struct {
		input  string
		output string
	}{
		{"", ""},
		{"\"", "\""},
		{"'", "'"},
		{"''", ""},
		{"\"\"", ""},
		{"'A'", "A"},
		{"\"A\"", "A"},
	}

	for _, d := range testdata {
		result := removeQoute(d.input)
		if result != d.output {
			t.Errorf("remove qoute expeted %s got %s", d.output, result)
		}
	}
}
