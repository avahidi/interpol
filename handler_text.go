package interpol

type TextHandler struct {
	text string
}

func NewTextHandler(text string, data *InterpolatorData) (Handler, error) {
	return &TextHandler{text: text}, nil
}

func (this *TextHandler) String() string {
	return this.text
}
func (this *TextHandler) Next() bool {
	return false
}

func (this *TextHandler) Reset() {
	// empty
}
