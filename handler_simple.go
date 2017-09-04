package interpol

import (
	"fmt"
)

// the text handler is a handler for static text
type textHandler struct {
	text string
}

func newTextHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	return &textHandler{text: text}, nil
}

func (t *textHandler) String() string {
	return t.text
}
func (t *textHandler) Next() bool {
	return false
}

func (t *textHandler) Reset() {
	// empty
}

// the copy handler will just copy another handlers value
type copyHandler struct {
	from Handler
}

func newCopyHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	from := ctx.tryImport(data.GetString("from", ""))
	if from == nil {
		return nil, fmt.Errorf("copy could not find target")
	}
	return &copyHandler{from: from}, nil
}

func (ch *copyHandler) String() string {
	return ch.from.String()
}
func (ch *copyHandler) Next() bool {
	return false
}

func (ch *copyHandler) Reset() {
	// empty
}

func init() {
	addDefaultHandlerFactory("copy", newCopyHandler)
}
