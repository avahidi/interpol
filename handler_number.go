package interpol

import (
	"fmt"
	"math/rand"
)

// counter handler is a up/down counter
type counterHandler struct {
	min, max, step int
	curr           int
	format         string
}

func newCounterHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
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

// random handler generates random in some interval
type randomHandler struct {
	start, width, count int
	curr, i             int
	format              string
}

func newrandomHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	min := data.GetInteger("min", 0)
	max := data.GetInteger("max", 100)
	if min >= max {
		return nil, fmt.Errorf("bad random min/max: %d/%d", min, max)
	}

	ret := &randomHandler{
		start:  min,
		width:  1 + max - min,
		count:  data.GetInteger("count", 5),
		format: data.GetString("format", "%d"),
	}

	if ret.count <= 0 {
		return nil, fmt.Errorf("bad random count: %d", ret.count)
	}

	ret.Reset()
	return ret, nil
}

func (rh *randomHandler) update() {
	rh.curr = rh.start + (rand.Int() % rh.width)
}

func (rh *randomHandler) String() string {
	return fmt.Sprintf(rh.format, rh.curr)
}
func (rh *randomHandler) Next() bool {
	if rh.i < rh.count {
		rh.i++
		rh.update()
	}
	return rh.i < rh.count
}

func (rh *randomHandler) Reset() {
	rh.i = 0
	rh.update()
}

func init() {
	addDefaultFactory("random", newrandomHandler)
}
