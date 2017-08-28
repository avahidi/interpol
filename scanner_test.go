package interpol

import (
	"testing"
)

func TestParseInterpol(t *testing.T) {
	id, err := parseInterpolator("justtype")
	if err != nil || id.Type != "justtype" || len(id.Properties) != 0 {
		t.Error("Failed to parse just interpol")
	}

	id, err = parseInterpolator("type1 data1 data2=value2")
	if err != nil || id.Type != "type1" || len(id.Properties) != 2 {
		t.Error("Failed to parse type1")
	}

	if id.Properties["data1"] != "" {
		t.Error("Failed to parse type1 data1")
	}

	if id.Properties["data2"] != "value2" {
		t.Error("Failed to parse type1 data2")
	}
}

func TestParseInterpolCorner(t *testing.T) {
	_, err := parseInterpolator("")
	if err == nil {
		t.Error("Failed to parse empty interpol")
	}

	id, err := parseInterpolator("spacetype datax  = valuex ")
	if err != nil || id.Properties["datax"] != "valuex" {
		t.Errorf("Failed to remove spaces")
	}
}

func TestParseLine(t *testing.T) {
	els, err := parseLine("justtext")
	if err != nil || len(els) != 1 || !els[0].isStatic || els[0].text != "justtext" {
		t.Error("failed to parse just text")
	}

	els, err = parseLine("{{justinterpol}}")
	if err != nil || len(els) != 1 || els[0].isStatic || els[0].text != "justinterpol" {
		t.Error("failed to parse just interpol")
	}

	els, err = parseLine("stuffbefore{{interpol}}stuffafter")
	if err != nil || len(els) != 3 {
		t.Error("failed to parse mixed")
	} else {
		if !els[0].isStatic || els[0].text != "stuffbefore" {
			t.Error("mixed part 1 failed to parse")
		}

		if els[1].isStatic || els[1].text != "interpol" {
			t.Error("mixed part 2 failed to parse")
		}
		if !els[2].isStatic || els[2].text != "stuffafter" {
			t.Error("mixed part 3 failed to parse")
		}
	}
}

func TestParseLineCorner(t *testing.T) {
	_, err := parseLine("")
	if err == nil {
		t.Error("failed to parse empty line")
	}
}
