package interpol

type textHandler struct {
	text string
}

func newTextHandler(text string, data *InterpolatorData) (Handler, error) {
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
