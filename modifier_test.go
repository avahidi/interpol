package interpol

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestSimpleModifier(t *testing.T) {
	var modifierTestData = []struct {
		modifier string
		indata   string
		outdata  string
	}{
		{"reverse", "", ""},
		{"reverse", "ABc", "cBA"},
		{"reverse", "頭いってる", "るてっい頭"},
		{"tolower", "Anders Jonas Ångström", "anders jonas ångström"},
		{"tolower", "Anders Jonas Ångström", "anders jonas ångström"},
		{"toupper", "Per Martin-Löf", "PER MARTIN-LÖF"},
		{"toupper", "Lars Bergström", "LARS BERGSTRÖM"},
		{"capitalize", "gordon GEKKO", "Gordon Gekko"},
		{"empty", "bla bla bla", ""},
		{"empty", "", ""},
		{"len", "five", "4"},
		{"base64", "", ""},
		{"base64", "Fermat said he had a proof", "RmVybWF0IHNhaWQgaGUgaGFkIGEgcHJvb2Y="},
		{"base64", "ずるい :)", "44Ga44KL44GEIDop"},

		{"trim", "", ""},
		{"trim", " AAA BB  C    \t\n", "AAA BB  C"},
	}

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

func TestBitflip(t *testing.T) {
	N := 1000
	indata := "Emacs is written in Lisp, which is the only computer language that is beautiful."
	for i := 0; i < N; i++ {
		ip := New()

		cmd := fmt.Sprintf("{{set data='%s' sep=$$ modifier=bitflip}}", indata)
		str, err := ip.Add(cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", cmd, err)
		} else {
			str2 := str.String()
			if len(str2) != len(indata) {
				t.Errorf("Bad size in bitflip")
			} else {
				cnt, diff := 0, 0
				for i := range indata {
					if indata[i] != str2[i] {
						cnt++
						diff = int(indata[i]) ^ int(str2[i])
					}
				}

				if cnt != 1 {
					t.Errorf("Not bytes were modified")
				}
				if (diff & (diff - 1)) != 0 { // diff is not a power of two?
					t.Errorf("Exactly one bit should have been changed")
				}
			}
		}
	}
}

func TestByteswap(t *testing.T) {
	N := 50
	indata := "0123456789" // no repeating data
	for i := 0; i < N; i++ {
		ip := New()

		cmd := fmt.Sprintf("{{set data='%s' sep=$$ modifier=byteswap}}", indata)
		str, err := ip.Add(cmd)
		if err != nil {
			t.Errorf("%s: failed to parse, %v", cmd, err)
		} else {
			str2 := str.String()
			if len(str2) != len(indata) {
				t.Errorf("Bad size in byteswap")
			} else {
				cnt, mask := 0, 0
				for i := range indata {
					if indata[i] != str2[i] {
						cnt++
					}
					// mask is a simple way to see if all bytes are still present
					if str2[i] >= '0' && str2[i] <= '9' {
						mask |= (1 << uint(str2[i]-'0'))
					}
				}
				if cnt != 2 || mask != (1<<10)-1 {
					t.Errorf("Exactly one byte should have been swapped")
				}
			}
		}
	}
}
