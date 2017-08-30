package test

import (
	"testing"

	"bitbucket.org/vahidi/interpol"
)

type interpolParserTestdata struct {
	log  string
	cmd  string
	typ  string
	data map[string]string
}

type lineParserTestdata struct {
	log      string
	line     string
	expected []interpol.TextElement
}

var interpolTestdata = []interpolParserTestdata{
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
	{"remove space 5", "spacetype nodata1 nodata2 =", "spacetype",
		map[string]string{"nodata1": "", "nodata2": ""}},		
	{"remove space 6", "spacetype nodata1 nodata2=", "spacetype",
		map[string]string{"nodata1": "", "nodata2": ""}},		
}

var lineTestdata = []lineParserTestdata{
	{"just text", "justtext", []interpol.TextElement{
		{Static: true, Text: "justtext"}}},
	{"just interpol", "{{justinterpol}}", []interpol.TextElement{
		{Static: false, Text: "justinterpol"}}},
	{"mixed line", "stuffbefore{{interpol}}stuffafter", []interpol.TextElement{
		{Static: true, Text: "stuffbefore"},
		{Static: false, Text: "interpol"},
		{Static: true, Text: "stuffafter"}}},
}

// test parsing of interpol commands
func TestParseInterpol(t *testing.T) {
	for _, test := range interpolTestdata {
		id, err := interpol.ParseInterpolator(test.cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", test.log, err)
		} else {
			if id.Type != test.typ {
				t.Errorf("%s: expected type %s got %s", test.log, test.typ, id.Type)
			}
			if len(id.Properties) != len(test.data) {
				t.Errorf("%s: expected %d properties, got %d",
					test.log, len(test.data), len(id.Properties))
			} else {
				for k, v := range test.data {
					if id.Properties[k] != v {
						t.Errorf("%s: expected %s=%s, got %s",
							test.log, k, v, id.Properties[k])
					}
				}
			}
		}
	}
}

func TestParseInterpolCorner(t *testing.T) {
	_, err := interpol.ParseInterpolator("")
	if err == nil {
		t.Error("should not allow empty interpol")
	}
}

// test parsing of lines:
func TestParseLine(t *testing.T) {
	for _, test := range lineTestdata {
		els, err := interpol.ParseLine(test.line)
		if err != nil {
			t.Errorf("%s: could not parse line, %v", test.log, err)
		} else {
			if len(test.expected) != len(els) {
				t.Errorf("%s: expected %d elements, got %d",
					test.log, len(test.expected), len(els))
			} else {
				for i, e := range test.expected {
					if e.Static != els[i].Static || e.Text != els[i].Text {
						t.Errorf("%s: element %d expected %v got %v",
							test.log, i, e, els[i])
					}
				}
			}
		}
	}
}

func TestParseLineCorner(t *testing.T) {
	_, err := interpol.ParseLine("")
	if err == nil {
		t.Error("should not allow empty line")
	}
}
