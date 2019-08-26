package interpol

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// element is part of an interpolated string
type element struct {
	handler  Handler
	modifier Modifier
}

func (ie *element) String() string {
	str := ie.handler.String()
	if ie.modifier != nil {
		str = ie.modifier.Modify(str)
	}
	return str
}

func (ie *element) Reset() {
	ie.handler.Reset()
}

func (ie *element) Next() bool {
	return ie.handler.Next()
}

// InterpolatedString contains an interpolator object
type InterpolatedString struct {
	elements []*element

	buffer *bytes.Buffer
}

func newInterpolatedString(size int) *InterpolatedString {
	ret := &InterpolatedString{}

	bs := make([]byte, 1024)
	ret.buffer = bytes.NewBuffer(bs)
	ret.elements = make([]*element, size)
	return ret
}

// convert InterpolatedString to a string
func (ips *InterpolatedString) String() string {
	ips.buffer.Reset()
	for _, e := range ips.elements {
		ips.buffer.WriteString(e.String())
	}
	return ips.buffer.String()
}

// Parameters for an interpolator command in a more accesible form
type Parameters struct {
	Type       string
	Properties map[string]string
}

// GetString returns an interpolation parameter as string
func (id *Parameters) GetString(name string, def string) string {
	if s, okay := id.Properties[name]; okay {
		return s
	}
	return def
}

// GetInteger returns an interpolation parameter as int
func (id *Parameters) GetInteger(name string, def int) int {
	if s, okay := id.Properties[name]; okay {
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
	}
	return def
}

// Interpol context for an interpolation
type Interpol struct {
	handlerFactories  map[string]HandlerFactory
	modifierFactories map[string]ModifierFactory
	elements          []*element
	exported          map[string]Handler
}

// New creates a new interpolator context
func New() *Interpol {
	ret := &Interpol{}
	ret.handlerFactories = make(map[string]HandlerFactory)
	ret.modifierFactories = make(map[string]ModifierFactory)
	ret.exported = make(map[string]Handler)
	ret.elements = make([]*element, 0)

	ret.Reset()
	return ret
}

// Reset resets everything to its original state
func (ip *Interpol) Reset() {
	for _, h := range ip.elements {
		h.Reset()
	}
}

// Next calculates the next value
func (ip *Interpol) Next() bool {
	for _, e := range ip.elements {
		if e.Next() {
			return true
		}
		e.Reset()
	}
	return false
}

//
// modifier functions
//

// AddModifier registers a new modifier
func (ip *Interpol) AddModifier(typ string, modifier ModifierFactory) error {
	if _, okay := ip.modifierFactories[typ]; okay {
		return fmt.Errorf("Modifier '%s' already exists", typ)
	}
	ip.modifierFactories[typ] = modifier
	return nil
}

func (ip *Interpol) findModifierFactory(name string) ModifierFactory {
	name = strings.ToLower(name)
	if def, okay := ip.modifierFactories[name]; okay {
		return def
	}
	return findDefaultModifierFactory(name)
}

func (ip *Interpol) createModifier(id *Parameters) (Modifier, error) {
	if id != nil {
		if name := id.GetString("modifier", ""); name != "" {
			mf := ip.findModifierFactory(name)
			if mf == nil {
				return nil, fmt.Errorf("Unknown modifier: %s", name)
			}
			return mf(ip, id)
		}
	}
	return nil, nil
}

//
// handler functions
//

// AddHandler adds a handler for a specific type of interpolator
func (ip *Interpol) AddHandler(typ string, creator HandlerFactory) error {
	if _, okay := ip.handlerFactories[typ]; okay {
		return fmt.Errorf("Handler for '%s' already exists", typ)
	}
	ip.handlerFactories[typ] = creator
	return nil
}

func (ip *Interpol) findHandlerFactory(name string) HandlerFactory {
	name = strings.ToLower(name)
	if def, okay := ip.handlerFactories[name]; okay {
		return def
	}
	return findDefaultHandlerFactory(name)
}

// import/export functions for copy

func (ip *Interpol) tryImport(name string) Handler {
	if h, found := ip.exported[name]; found {
		return h
	}
	return nil
}

func (ip *Interpol) tryExport(data *Parameters, h Handler) error {
	if data != nil {
		if name, okay := data.Properties["name"]; okay {
			if _, seenbefore := ip.exported[name]; seenbefore {
				return fmt.Errorf("name '%s' already exists", name)
			}
			ip.exported[name] = h
		}
	}
	return nil
}

// Add creates a new string to be interpolated
func (ip *Interpol) Add(text string) (*InterpolatedString, error) {

	// parse the line for elements
	els, err := parseLine(text)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse '%s'", text)
	}

	// convert elements to handlers and build the return string
	ret := newInterpolatedString(len(els))
	for i, e := range els {
		var factory HandlerFactory
		var id *Parameters
		if e.static {
			factory = newTextHandler
			id = nil
		} else {
			id, err = parseInterpolator(e.text)
			if err != nil {
				return nil, err
			}
			factory = ip.findHandlerFactory(id.Type)
			if factory == nil {
				return nil, fmt.Errorf("Cannot find a handler for '%s'", e.text)
			}
		}

		handler, err := factory(ip, e.text, id)
		if err != nil {
			return nil, fmt.Errorf("Cannot initialize handler '%s': %v", text, err)
		}

		modifier, err := ip.createModifier(id)
		if err != nil {
			return nil, fmt.Errorf("Cannot create modifier '%s': %v", text, err)
		}

		err = ip.tryExport(id, handler)
		if err != nil {
			return nil, err
		}

		ret.elements[i] = &element{handler: handler, modifier: modifier}

	}

	// add the new handlers to the list of all handlers in this context
	for _, h := range ret.elements {
		ip.elements = append(ip.elements, h)
	}
	return ret, nil
}

// AddMultiple creates multiple strings to be interpolated
func (ip *Interpol) AddMultiple(texts ...string) ([]*InterpolatedString, error) {
	ret := make([]*InterpolatedString, len(texts))
	for i, text := range texts {
		ips, err := ip.Add(text)
		if err != nil {
			return nil, err
		}
		ret[i] = ips
	}
	return ret, nil
}
