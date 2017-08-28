package interpol

import (
	"fmt"
	"math/rand"
)

type RandomHandler struct {
	start, width, count int
	curr, i             int
	format              string
}

func NewRandomHandler(text string, data *InterpolatorData) (Handler, error) {
	min := data.GetInteger("min", 0)
	max := data.GetInteger("max", 100)
	if min >= max {
		return nil, fmt.Errorf("Bad random min/max: %d/%d\n", min, max)
	}

	ret := &RandomHandler{
		start:  min,
		width:  1 + max - min,
		count:  data.GetInteger("count", 5),
		format: data.GetString("format", "%d"),
	}

	if ret.count <= 0 {
		return nil, fmt.Errorf("Bad random count: %d\n", ret.count)
	}

	ret.Reset()
	return ret, nil
}

func (this *RandomHandler) update() {
	this.curr = this.start + (rand.Int() % this.width)
}

func (this *RandomHandler) String() string {
	return fmt.Sprintf(this.format, this.curr)
}
func (this *RandomHandler) Next() bool {
	if this.i < this.count {
		this.i++
		this.update()
	}
	return this.i < this.count
}

func (this *RandomHandler) Reset() {
	this.i = 0
	this.update()
}
