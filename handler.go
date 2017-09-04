package interpol

// Handler represents a handler for a certain type of interpolation
type Handler interface {
	String() string
	Next() bool
	Reset()
}

// HandlerFactory creates a new handler for a given text or command
type HandlerFactory func(ctx *Interpol, text string, data *InterpolatorData) (Handler, error)

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
