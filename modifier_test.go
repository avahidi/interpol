package interpol

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

type modifierTestDataFormat struct {
	modifier string
	indata   string
	outdata  string
}

var modifierTestData = []modifierTestDataFormat{
	{"tolower", "Anders Jonas Ångström", "anders jonas ångström"},
	{"tolower", "Anders Jonas Ångström", "anders jonas ångström"},
	{"toupper", "Per Martin-Löf", "PER MARTIN-LÖF"},
	{"toupper", "Lars Bergström", "LARS BERGSTRÖM"},
	{"capitalize", "gordon GEKKO", "Gordon Gekko"},
}

func TestSimpleModifier(t *testing.T) {
	for _, test := range modifierTestData {
		cmd := fmt.Sprintf("{{set data='%s' sep='%s' modifier=%s}}",
			test.indata, "$$$", test.modifier)
		ip := New()

		str, err := ip.Add(cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", cmd, err)
		} else {
			if str.String() != test.outdata {
				t.Errorf("expected %s got %s", str, test.outdata)
			}
		}
	}
}

// since Leet is random we cant test it the normal way,
// instead we are looking at average of N runs
func TestLeetModifier(t *testing.T) {
	N := 1000
	indata := "I do my taxes in emacs..."
	counters := make([]int, utf8.RuneCountInString(indata))

	for i := 0; i < N; i++ {
		ip := New()
		cmd := fmt.Sprintf("{{set data='%s' sep=$$ modifier=1337}}", indata)
		str, err := ip.Add(cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", cmd, err)
		} else {
			str2 := str.String()
			if strings.ToLower(str2) != strings.ToLower(indata) {
				t.Errorf("'%s' is not leet speak for '%s'", str, indata)
			}
			for i, c := range str2 {
				if unicode.IsUpper(c) {
					counters[i]++
				}
			}
		}
	}

	// really bad statistical test:
	for i, c := range indata {
		if unicode.IsLetter(c) {
			n := counters[i]
			if 100*n < (15*N) || 100*n > (85*N) {
				t.Errorf("Leet modifier charcter %d upper rate was %d%%", i, 100*n/N)
			}
		}
	}

}
