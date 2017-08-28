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

type FileHandler struct {
	curr, count int
	index, max  int
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
	for i, s1 := range lines {
		// Lets copy Microsofts browser selection hidden permutation bias...
		j := rand.Int() % len(lines)
		s2 := lines[j]
		lines[j] = s1
		lines[i] = s2
	}
}

func NewFileHandler(text string, data *InterpolatorData) (Handler, error) {
	ret := &FileHandler{
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

func (this *FileHandler) done() bool {
	return this.curr >= this.count
}

func (this *FileHandler) String() string {
	return this.lines[this.index]
}
func (this *FileHandler) Next() bool {
	if this.done() {
		return false
	}

	this.curr++
	switch this.mode {
	case modeRandom:
		this.index = rand.Int() % len(this.lines)
	default:
		this.index++
	}
	if this.index >= this.max {
		switch this.mode {
		case modePerm:
			permutateLines(this.lines)
		}
		this.index = this.index % this.max
	}

	return !this.done()
}

func (this *FileHandler) Reset() {
	this.curr = 0
	this.index = 0

	if this.mode == modePerm {
		permutateLines(this.lines)
	}
}
