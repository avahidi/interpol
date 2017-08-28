package interpol

import (
	"fmt"
)

type CounterHandler struct {
	min, max, step int
	curr           int
	format         string
}

func NewCounterHandler(text string, data *InterpolatorData) (Handler, error) {
	ret := &CounterHandler{
		min:    data.GetInteger("min", 0),
		max:    data.GetInteger("max", 10),
		step:   data.GetInteger("step", 1),
		format: data.GetString("format", "%d"),
	}
	ret.Reset()
	return ret, nil
}

func (this *CounterHandler) done() bool {
	return (this.step > 0 && this.curr > this.max) || (this.step < 0 && this.curr < this.max)
}

func (this *CounterHandler) String() string {
	return fmt.Sprintf(this.format, this.curr)
}
func (this *CounterHandler) Next() bool {
	if !this.done() {
		this.curr += this.step
	}
	return !this.done()
}

func (this *CounterHandler) Reset() {
	this.curr = this.min
}
