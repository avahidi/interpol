package interpol

import (
	"math/rand"
	"strings"
	"unicode"
)

// Modifier interface represents any functions that modifies an
// interpolator
type Modifier interface {
	Modify(string) string
}

// ModifierFactory creates a new modifier
type ModifierFactory func(ctx *Interpol, data *InterpolatorData) (Modifier, error)

//
var defaultModifierFactories = map[string]ModifierFactory{
	"toupper":    newToupperModifier,
	"tolower":    newTolowerModifier,
	"capitalize": newCapitalizeModifier,
	"leet":       newLeetModifier,
	"1337":       newLeetModifier,
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
// to upper
//

type toupperModifier struct{}

func (t *toupperModifier) Modify(str string) string { return strings.ToUpper(str) }

func newToupperModifier(ctx *Interpol, data *InterpolatorData) (Modifier, error) {
	return &toupperModifier{}, nil
}

//
// to lower
//

type tolowerModifier struct{}

func (t *tolowerModifier) Modify(str string) string { return strings.ToLower(str) }

func newTolowerModifier(ctx *Interpol, data *InterpolatorData) (Modifier, error) {
	return &tolowerModifier{}, nil
}

//
// to lower
//

type capitalizeModifier struct{}

func (c *capitalizeModifier) Modify(str string) string {
	return strings.Title(strings.ToLower(str))
}

func newCapitalizeModifier(ctx *Interpol, data *InterpolatorData) (Modifier, error) {
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

func newLeetModifier(ctx *Interpol, data *InterpolatorData) (Modifier, error) {
	return &leetModifier{}, nil
}
