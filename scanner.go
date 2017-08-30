package interpol

import (
	"bytes"
	"fmt"
	"strings"
)

// TextElement represents a sub-string that is either static or an interpolator
// (this internal data structure is public mainly to simplify testing)
type TextElement struct {
	Static bool
	Text   string
}


// remove spaces when the command looks like this
// cmd prop = value
// cmd prop1 prop2 = value
// cmd prop1 prop2 =value
func removeSpacesInCommand(list []string) [] string {
	ret := make([]string, 0)
	ret = append(ret, list[0])
	
	for i, last := 1, len(list)-1; i <= last; i++ {
		out := list[i]	
		if strings.Index(list[i], "=") == - 1 && i < last {			
			if next := strings.Trim(list[i + 1], " \t") ; next[0] == '=' {				
				out = out + next
				i++
				if next == "=" && i < last && strings.Index(list[i + 1], "=") == -1 {
					out = out + strings.Trim(list[i + 1], " \t") 
					i++						
				}
			}

		}
		ret = append(ret, out)
	}

	return ret
}

// ParseInterpolator parses the interpolator TextElements.
// It is basically a betterversion of strings.Split() that handles ' and " and \
func ParseInterpolator(text string) (*InterpolatorData, error) {
	const (
		normal = iota
		waitMeta
		waitSingle
		waitDouble
	)
	b := bytes.NewBuffer(nil)
	state := normal
	list := make([]string, 0)

	for _, r := range text {
		switch state {
		case normal:
			switch r {
			case '\'':
				state = waitSingle
			case '"':
				state = waitDouble
			case '\\':
				state = waitMeta
			case ' ', ',':
				if b.Len() > 0 {
					list = append(list, b.String())
					b.Reset()
				}
			default:
				if b.Len() > 0 || (r != ' ' && r != '\t') {
					b.WriteRune(r)
				}
			}
		case waitSingle:
			if r == '\'' {
				state = normal
			} else {
				b.WriteRune(r)
			}
		case waitDouble:
			if r == '"' {
				state = normal
			} else {
				b.WriteRune(r)
			}
		case waitMeta:
			if r != '\\' {
				b.WriteRune(r)
			}
			state = normal

		}
	}

	if b.Len() > 0 {
		list = append(list, b.String())
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("interpolator contains no data")
	}

	if strings.Index(list[0], "=") != -1 {
		return nil, fmt.Errorf("type invalid: %s", text)
	}

	list = removeSpacesInCommand(list)

	i := &InterpolatorData{Type: list[0], Properties: make(map[string]string)}
	for _, p := range list[1:] {
		p0 := strings.Trim(p, " \t")
		m := strings.Index(p0, "=")
		if m == -1 {
			i.Properties[p0] = ""
		} else {
			pl := strings.Trim(p0[:m], " \t")
			pr := strings.Trim(p0[m+1:], " \t")
			i.Properties[pl] = pr
		}
	}

	return i, nil
}

// ParseLine divides a line into a number of TextElements that
// are either a static string or an interpolator description
func ParseLine(line string) ([]TextElement, error) {

	if len(line) == 0 {
		return []TextElement{}, fmt.Errorf("Empty line")
	}

	ret := make([]TextElement, 0)
	for len(line) > 0 {
		n := strings.Index(line, "{{")
		if n == -1 {
			ret = append(ret, TextElement{Static: true, Text: line})
			line = ""
		} else {
			if n != 0 {
				ret = append(ret, TextElement{Static: true, Text: line[:n]})
			}
			line = line[n+2:]
			m := strings.Index(line, "}}")
			if m == -1 {
				ret = append(ret, TextElement{Static: false, Text: line})
				line = ""
			} else {
				ret = append(ret, TextElement{Static: false, Text: line[:m]})
				line = line[m+2:]
			}
		}
	}
	return ret, nil
}
