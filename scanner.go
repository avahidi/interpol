package interpol

import (
	"bytes"
	"fmt"
	"strings"
)

type element struct {
	isStatic bool
	text     string
}

// like strings.Split() but handles ' and " and \
func parseInterpolator(text string) (*InterpolatorData, error) {
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

func parseLine(line string) ([]element, error) {

	if len(line) == 0 {
		return []element{}, fmt.Errorf("Empty line")
	}

	ret := make([]element, 0)
	for len(line) > 0 {
		n := strings.Index(line, "{{")
		if n == -1 {
			ret = append(ret, element{isStatic: true, text: line})
			line = ""
		} else {
			if n != 0 {
				ret = append(ret, element{isStatic: true, text: line[:n]})
			}
			line = line[n+2:]
			m := strings.Index(line, "}}")
			if m == -1 {
				ret = append(ret, element{isStatic: false, text: line})
				line = ""
			} else {
				ret = append(ret, element{isStatic: false, text: line[:m]})
				line = line[m+2:]
			}
		}
	}
	return ret, nil
}
