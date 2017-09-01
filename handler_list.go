package interpol

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// processing mode
const (
	modeLinear = iota
	modePerm
	modeRandom
)

var modeMap = map[string]int{
	"linear":      modeLinear,
	"perm":        modePerm,
	"permutation": modePerm,
	"random":      modeRandom,
	"rand":        modeRandom,
}

type listHandler struct {
	curr, count int
	index, max  int // for items
	mode        int
	items       []string
}

// some helper functions
func readFileitems(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]string, 0)
	rd := bufio.NewReader(file)

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		str = strings.Trim(str, " \t\n\r")
		if len(str) > 0 {
			ret = append(ret, str)
		}
	}

	return ret, nil
}

func permutateitems(items []string) {
	if len(items) < 2 {
		return
	}
	// first loop ensures that we never repeat ourselves
	last := items[len(items)-1]
	for {
		// permutation bias? well it worked for Microsofts browser selection...
		for i := range items {
			j := rand.Int() % len(items)
			items[i], items[j] = items[j], items[i]
		}

		if last != items[0] {
			return
		}
	}
}

func newFileHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	// get file contents
	filename := data.GetString("filename", "")
	if filename == "" {
		return nil, fmt.Errorf("no filename was given")
	}

	items, err := readFileitems(filename)
	if err != nil {
		return nil, err
	}

	return newListHandler(items, data)
}

func newSetHandler(ctx *Interpol, text string, data *InterpolatorData) (Handler, error) {
	var items []string
	sep := data.GetString("sep", "")
	set := data.GetString("data", "")

	items = strings.Split(set, sep)
	return newListHandler(items, data)
}

func newListHandler(items []string, data *InterpolatorData) (Handler, error) {
	ret := &listHandler{
		count: data.GetInteger("count", -1),
		items: items,
		max:   len(items),
	}

	// max items and output count
	if ret.max == 0 {
		return nil, fmt.Errorf("Empty file")
	}

	// user didn't specify count...
	if ret.count <= 0 {
		ret.count = ret.max
	}

	// get processing mode
	modename := data.GetString("mode", "linear")
	mode, okay := modeMap[modename]
	if !okay {
		return nil, fmt.Errorf("unknown file mode")
	}
	ret.mode = mode

	ret.Reset()
	return ret, nil
}

func (fh *listHandler) done() bool {
	return fh.curr >= fh.count
}

func (fh *listHandler) String() string {
	return fh.items[fh.index]
}

func (fh *listHandler) Next() bool {
	if fh.done() {
		return false
	}

	fh.curr++
	switch fh.mode {
	case modeRandom:
		fh.index = rand.Int() % len(fh.items)
	default:
		fh.index++
		if fh.index >= fh.max {
			if fh.mode == modePerm {
				permutateitems(fh.items)
			}
			fh.index = fh.index % fh.max
		}
	}

	return !fh.done()
}

func (fh *listHandler) Reset() {
	fh.curr = 0
	fh.index = 0

	if fh.mode == modePerm {
		permutateitems(fh.items)
	}
}

func init() {
	addDefaultFactory("file", newFileHandler)
	addDefaultFactory("set", newSetHandler)
}
