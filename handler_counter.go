package interpol

import (
	"fmt"
)

type counterHandler struct {
	min, max, step int
	curr           int
	format         string
}

func newCounterHandler(text string, data *InterpolatorData) (Handler, error) {
	ret := &counterHandler{
		min:    data.GetInteger("min", 0),
		max:    data.GetInteger("max", 10),
		step:   data.GetInteger("step", 1),
		format: data.GetString("format", "%d"),
	}
	ret.Reset()
	return ret, nil
}

func (ch *counterHandler) done() bool {
	return (ch.step > 0 && ch.curr > ch.max) || (ch.step < 0 && ch.curr < ch.max)
}

func (ch *counterHandler) String() string {
	return fmt.Sprintf(ch.format, ch.curr)
}
func (ch *counterHandler) Next() bool {
	if !ch.done() {
		ch.curr += ch.step
	}
	return !ch.done()
}

func (ch *counterHandler) Reset() {
	ch.curr = ch.min
}

func init() {
	addDefaultFactory("counter", newCounterHandler)
}
