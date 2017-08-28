package interpol

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var baseHandlerFactories = map[string]HandlerFactory{
	"text":    NewTextHandler,
	"counter": NewCounterHandler,
	"random":  NewRandomHandler,
	"file":    NewFileHandler,
}

// InterpolatedString
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

func (this *InterpolatedString) String() string {
	this.buffer.Reset()
	for _, h := range this.handlers {
		this.buffer.WriteString(h.String())
	}
	return this.buffer.String()
}

// InterpolatorData
type InterpolatorData struct {
	Type       string
	Properties map[string]string
}

func (this *InterpolatorData) GetString(name string, def string) string {
	if s, okay := this.Properties[name]; okay {
		return s
	}
	return def
}
func (this *InterpolatorData) GetInteger(name string, def int) int {
	if s, okay := this.Properties[name]; okay {
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

type HandlerFactory func(text string, data *InterpolatorData) (Handler, error)

// Interpol contains the context for everything
type Interpol struct {
	factory  map[string]HandlerFactory
	handlers []Handler
}

func NewInterpol() *Interpol {
	ret := &Interpol{}
	ret.factory = make(map[string]HandlerFactory)
	ret.handlers = make([]Handler, 0)

	// register the base handlers
	for k, v := range baseHandlerFactories {
		ret.AddHandler(k, v)
	}
	return ret
}

func (this *Interpol) Reset() {
	for _, h := range this.handlers {
		h.Reset()
	}
}

func (this *Interpol) Next() bool {
	for i := 0; i < len(this.handlers); i++ {
		if this.handlers[i].Next() {
			return true
		}
		this.handlers[i].Reset()
	}
	return false
}

func (this *Interpol) AddHandler(type_ string, creator HandlerFactory) error {
	if _, okay := this.factory[type_]; okay {
		return fmt.Errorf("Handler for '%s' already exists", type_)
	}
	this.factory[type_] = creator
	return nil
}

func (this *Interpol) Add(text string) (*InterpolatedString, error) {

	// parse the line for elements
	els, err := parseLine(text)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse '%s'", text)
	}

	// convert elements to handlers and build the return string
	ret := newInterpolatedString(len(els))
	for i, e := range els {
		var factory HandlerFactory
		var id *InterpolatorData
		if e.isStatic {
			factory = NewTextHandler
			id = nil
		} else {
			id, err = parseInterpolator(e.text)
			if err != nil {
				return nil, err
			}
			var okay bool
			factory, okay = this.factory[strings.ToLower(id.Type)]
			if !okay {
				return nil, fmt.Errorf("Cannot find a handler for '%s'", e.text)
			}
		}

		handler, err := factory(e.text, id)
		if err != nil {
			return nil, fmt.Errorf("Cannot initialize handler '%s': %v", text, err)
		}

		ret.handlers[i] = handler
	}

	// add the new handlers to the list of all handlers in this context
	for _, h := range ret.handlers {
		this.handlers = append(this.handlers, h)
	}
	return ret, nil
}
