package interpol

import (
	"fmt"
	"math/rand"
)

type randomHandler struct {
	start, width, count int
	curr, i             int
	format              string
}

func newrandomHandler(text string, data *InterpolatorData) (Handler, error) {
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
