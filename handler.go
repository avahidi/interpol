package interpol

// Handler represents a handler for a certain type of interpolation
//
// Note 1: Handler object should call Reset() when created
// Note 2: String() should be valid right after Reset()
// Note 3: Unlike Interpol, Next() after Reset() yields the _second_ value not first.
type Handler interface {
	String() string

	// Get the next value, returns false if no more values are available
	Next() bool

	// Reset the handler and prepare the first value
	Reset()
}

// HandlerFactory creates a new handler for a given text or command
type HandlerFactory func(ctx *Interpol, text string, data *Parameters) (Handler, error)

// default factories are added to all interpolators upon creation
var defaultHandlerFactories = map[string]HandlerFactory{}

func addDefaultHandlerFactory(name string, factory HandlerFactory) {
	defaultHandlerFactories[name] = factory
}

func findDefaultHandlerFactory(name string) HandlerFactory {
	if fact, okay := defaultHandlerFactories[name]; okay {
		return fact
	}
	return nil
}
