package interpol

import (
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

// Modifier interface represents any functions that modifies an
// interpolator
type Modifier interface {
	Modify(string) string
}

// ModifierFactory creates a new modifier
type ModifierFactory func(ctx *Interpol, data *Parameters) (Modifier, error)

//
var defaultModifierFactories = map[string]ModifierFactory{
	"reverse":    newReverseModifier,
	"trim":       newTrimModifier,
	"toupper":    newToupperModifier,
	"tolower":    newTolowerModifier,
	"capitalize": newCapitalizeModifier,
	"leet":       newLeetModifier,
	"1337":       newLeetModifier,
	"empty":      newEmptyModifier,
	"len":        newLenModifier,
	"bitflip":    newBitflipModifier,
	"byteswap":   newByteswapModifier,
	"base64":     newBase64Modifier,
}

func addDefaultModifierFactories(name string, factory ModifierFactory) {
	defaultModifierFactories[name] = factory
}

func findDefaultModifierFactory(name string) ModifierFactory {
	if fact, okay := defaultModifierFactories[name]; okay {
		return fact
	}
	return nil
}

//
// reverse
//

type reverseModifier struct{}

func (t *reverseModifier) Modify(str string) string {
	n := len(str)
	rr := make([]rune, n)
	for _, v := range str {
		n--
		rr[n] = v
	}
	return string(rr[n:])
}

func newReverseModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &reverseModifier{}, nil
}

//
// trim
//

type trimModifier struct{}

func (t *trimModifier) Modify(str string) string { return strings.Trim(str, " \t\n") }

func newTrimModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &trimModifier{}, nil
}

//
// to upper
//

type toupperModifier struct{}

func (t *toupperModifier) Modify(str string) string { return strings.ToUpper(str) }

func newToupperModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &toupperModifier{}, nil
}

//
// to lower
//

type tolowerModifier struct{}

func (t *tolowerModifier) Modify(str string) string { return strings.ToLower(str) }

func newTolowerModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &tolowerModifier{}, nil
}

//
// capitalize
//

type capitalizeModifier struct{}

func (c *capitalizeModifier) Modify(str string) string {
	return strings.Title(strings.ToLower(str))
}

func newCapitalizeModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &capitalizeModifier{}, nil
}

//
// leet speak
//

type leetModifier struct{}

func (t *leetModifier) Modify(str string) string {
	return strings.Map(func(c rune) rune {
		if (rand.Int() & 1) == 0 {
			return unicode.ToUpper(c)
		}
		return unicode.ToLower(c)
	}, str)
}

func newLeetModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &leetModifier{}, nil
}

//
// empty
//

type emptyModifier struct{}

func (t *emptyModifier) Modify(str string) string {
	return ""
}

func newEmptyModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &emptyModifier{}, nil
}

//
// len
//

type lenModifier struct{}

func (t *lenModifier) Modify(str string) string {
	return strconv.Itoa(len(str))
	// return strconv.Itoa(utf8.RuneCountInString(str))
}

func newLenModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &lenModifier{}, nil
}

//
// bitflip
//

type bitflipModifier struct{}

func (t *bitflipModifier) Modify(str string) string {
	if str == "" {
		return str
	}

	// probably not the most efficient way, but this is what we got
	i, b := rand.Int()%len(str), rand.Uint32()%8
	bs := []byte(str)
	bs[i] = bs[i] ^ (1 << b)
	return string(bs)
}

func newBitflipModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &bitflipModifier{}, nil
}

//
// byteswap
//

type byteswapModifier struct{}

func (t *byteswapModifier) Modify(str string) string {
	if len(str) < 2 {
		return str
	}

	bs := []byte(str)
	p1, p2 := rand.Int()%len(str), rand.Int()%len(str)

	// XXX: this could be made more efficient
	for p1 == p2 {
		p2 = rand.Int() % len(str)
	}

	bs[p1], bs[p2] = bs[p2], bs[p1]
	return string(bs)
}

func newByteswapModifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &byteswapModifier{}, nil
}

// base64

type base64Modifier struct{}

func (t *base64Modifier) Modify(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func newBase64Modifier(ctx *Interpol, data *Parameters) (Modifier, error) {
	return &base64Modifier{}, nil
}
