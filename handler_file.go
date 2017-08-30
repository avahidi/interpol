package interpol

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// file processing mode
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

type fileHandler struct {
	curr, count int
	index, max  int // for lines
	mode        int
	lines       []string
}

// some helper functions
func readFileLines(filename string) ([]string, error) {
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

func permutateLines(lines []string) {
	if len(lines) < 2 {
		return
	}
	// first loop ensures that we never repeat ourselves
	last := lines[len(lines)-1]
	for {
		// permutation bias? well it worked for Microsofts browser selection...
		for i := range lines {
			j := rand.Int() % len(lines)
			lines[i], lines[j] = lines[j], lines[i]
		}

		if last != lines[0] {
			return
		}
	}

}

func newfileHandler(text string, data *InterpolatorData) (Handler, error) {
	ret := &fileHandler{
		count: data.GetInteger("count", -1),
	}

	// get file contents
	filename := data.GetString("filename", "")
	if filename == "" {
		return nil, fmt.Errorf("no filename was given")
	}

	lines, err := readFileLines(filename)
	if err != nil {
		return nil, err
	}
	ret.lines = lines

	// max lines and output count
	ret.max = len(ret.lines)
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

func (fh *fileHandler) done() bool {
	return fh.curr >= fh.count
}

func (fh *fileHandler) String() string {
	return fh.lines[fh.index]
}

func (fh *fileHandler) Next() bool {
	if fh.done() {
		return false
	}

	fh.curr++
	switch fh.mode {
	case modeRandom:
		fh.index = rand.Int() % len(fh.lines)
	default:
		fh.index++
		if fh.index >= fh.max {
			if fh.mode == modePerm {
				permutateLines(fh.lines)
			}
			fh.index = fh.index % fh.max
		}
	}

	return !fh.done()
}

func (fh *fileHandler) Reset() {
	fh.curr = 0
	fh.index = 0

	if fh.mode == modePerm {
		permutateLines(fh.lines)
	}
}

func init() {
	addDefaultFactory("file", newfileHandler)
}
