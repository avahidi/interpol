package interpol

import (
	"strings"
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
	"toupper": newToupperModifier,
	"tolower": newTolowerModifier,
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
