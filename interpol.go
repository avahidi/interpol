package interpol

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// default factories are added to all interpolators upon creation
var defaultFactories = map[string]HandlerFactory{}

func addDefaultFactory(name string, factory HandlerFactory) {
	defaultFactories[name] = factory
}

// InterpolatedString contains an interpolator object
type InterpolatedString struct {
	handlers []Handler
	buffer   *bytes.Buffer
}

func newInterpolatedString(size int) *InterpolatedString {
	ret := &InterpolatedString{}

	bs := make([]byte, 1024)
	ret.buffer = bytes.NewBuffer(bs)
	ret.handlers = make([]Handler, size)
	return ret
}

// convert InterpolatedString to a string
func (ips *InterpolatedString) String() string {
	ips.buffer.Reset()
	for _, h := range ips.handlers {
		ips.buffer.WriteString(h.String())
	}
	return ips.buffer.String()
}

// InterpolatorData interpolator command in a more accesible form
type InterpolatorData struct {
	Type       string
	Properties map[string]string
}

// GetString extracts a string from interpolation data
func (id *InterpolatorData) GetString(name string, def string) string {
	if s, okay := id.Properties[name]; okay {
		return s
	}
	return def
}

// GetInteger extracts an integer from interpolation data
func (id *InterpolatorData) GetInteger(name string, def int) int {
	if s, okay := id.Properties[name]; okay {
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
	}
	return def
}

// Handler represents a handler for a certain type of interpolation
type Handler interface {
	String() string
	Next() bool
	Reset()
}

// HandlerFactory creates a new handler for a given text or command
type HandlerFactory func(text string, data *InterpolatorData) (Handler, error)

// Interpol context for an interpolation
type Interpol struct {
	factory  map[string]HandlerFactory
	handlers []Handler
}

// New creates a new interpolator context
func New() *Interpol {
	ret := &Interpol{}
	ret.factory = make(map[string]HandlerFactory)
	ret.handlers = make([]Handler, 0)

	// register the base handlers
	for k, v := range defaultFactories {
		ret.AddHandler(k, v)
	}
	return ret
}

// Reset resets everything to its original state
func (ip *Interpol) Reset() {
	for _, h := range ip.handlers {
		h.Reset()
	}
}

// Next calculates the next value
func (ip *Interpol) Next() bool {
	for i := 0; i < len(ip.handlers); i++ {
		if ip.handlers[i].Next() {
			return true
		}
		ip.handlers[i].Reset()
	}
	return false
}

// AddHandler adds a handler for a specific type of interpolator
func (ip *Interpol) AddHandler(typ string, creator HandlerFactory) error {
	if _, okay := ip.factory[typ]; okay {
		return fmt.Errorf("Handler for '%s' already exists", typ)
	}
	ip.factory[typ] = creator
	return nil
}

// Add creates a new string to be interpolated
func (ip *Interpol) Add(text string) (*InterpolatedString, error) {

	// parse the line for elements
	els, err := ParseLine(text)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse '%s'", text)
	}

	// convert elements to handlers and build the return string
	ret := newInterpolatedString(len(els))
	for i, e := range els {
		var factory HandlerFactory
		var id *InterpolatorData
		if e.Static {
			factory = newTextHandler
			id = nil
		} else {
			id, err = ParseInterpolator(e.Text)
			if err != nil {
				return nil, err
			}
			var okay bool
			factory, okay = ip.factory[strings.ToLower(id.Type)]
			if !okay {
				return nil, fmt.Errorf("Cannot find a handler for '%s'", e.Text)
			}
		}

		handler, err := factory(e.Text, id)
		if err != nil {
			return nil, fmt.Errorf("Cannot initialize handler '%s': %v", text, err)
		}

		ret.handlers[i] = handler
	}

	// add the new handlers to the list of all handlers in this context
	for _, h := range ret.handlers {
		ip.handlers = append(ip.handlers, h)
	}
	return ret, nil
}
