package interpol

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type scantype int

const (
	EOF scantype = iota
	Text
	Operator
	Error
)

type scanner struct {
	input string
	curr  int
}

func (l *scanner) next() (string, scantype) {
	//
	for l.curr < len(l.input) {
		rune, width := utf8.DecodeRuneInString(l.input[l.curr:])
		if rune != ' ' && rune != '\t' && rune != ',' {
			break
		}
		l.curr += width
	}

	start, qoute, meta := l.curr, false, false
	for l.curr < len(l.input) {
		rune, width := utf8.DecodeRuneInString(l.input[l.curr:])
		if meta {
			meta = false
		} else if qoute {
			if rune == '\'' || rune == '"' {
				qoute = false
			}
		} else {
			if rune == '\\' {
				meta = true
			} else if rune == '\'' || rune == '"' {
				qoute = true
			} else if rune == ' ' || rune == '\t' || rune == ',' || (rune == '=' && l.curr != start) {
				break
			} else if rune == '=' {
				l.curr += width
				break
			}
		}
		l.curr += width
	}

	str := l.input[start:l.curr]
	// fmt.Printf("START: %d, curr=%d str='%s'\n", start, l.curr, str)
	if start == l.curr {
		return "", EOF
	}

	// unexpected end
	if qoute || meta {
		return "", Error
	}

	if str == "=" {
		return str, Operator
	}

	// remove qoute marks if any
	if len(str) > 1 && ((strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'")) ||
		(strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\""))) {
		str = str[1 : len(str)-1]
	}
	return str, Text

}

func newScanner(str string) *scanner {
	return &scanner{input: str}
}

// parseInterpolator parses the interpolator textElements.
// It is basically a betterversion of strings.Split() that handles ' and " and \
func parseInterpolator(text string) (*InterpolatorData, error) {
	ret := &InterpolatorData{
		Properties: make(map[string]string),
	}

	s := newScanner(text)

	state := 0
	varname := ""
	for {

		str, typ := s.next()
		if typ == EOF {
			if state != 1 {
				return nil, fmt.Errorf("Unexpected end in '%s'", text)
			}
			break
		}

		switch state {
		case 0:
			if typ != Text {
				return nil, fmt.Errorf("Expected type, got %s", str)
			}
			ret.Type = str
			state = 1
		case 1:
			if typ == Operator {
				state = 2
			} else {
				varname = str
				ret.Properties[str] = ""
			}
		case 2:
			if varname == "" {
				return nil, fmt.Errorf("Unexpected '=' in '%s'", text)
			}
			ret.Properties[varname] = str
			varname = ""
			state = 1
		}
	}
	return ret, nil
}

// textElement represents a sub-string that is either static or an interpolator
type textElement struct {
	static bool
	text   string
}

// ParseLine divides a line into a number of textElements that
// are either a static string or an interpolator description
func parseLine(line string) ([]textElement, error) {

	if len(line) == 0 {
		return []textElement{}, fmt.Errorf("Empty line")
	}

	ret := make([]textElement, 0)
	for len(line) > 0 {
		n := strings.Index(line, "{{")
		if n == -1 {
			ret = append(ret, textElement{static: true, text: line})
			line = ""
		} else {
			if n != 0 {
				ret = append(ret, textElement{static: true, text: line[:n]})
			}
			line = line[n+2:]
			m := strings.Index(line, "}}")
			if m == -1 {
				ret = append(ret, textElement{static: false, text: line})
				line = ""
			} else {
				ret = append(ret, textElement{static: false, text: line[:m]})
				line = line[m+2:]
			}
		}
	}
	return ret, nil
}
