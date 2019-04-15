package interpol

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadFile is a helper function to read non-empty lines from a file.
// The reason this is a public function is because it is also used by Police
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {

		// this is a special case for snaps that have no access to $HOME
		// if err == os.ErrPermission {
		if strings.Contains(err.Error(), os.ErrPermission.Error()) {
			fmt.Fprintf(os.Stderr, "\n\n"+
				"** Could not open file due to permission issues.                 **\n"+
				"** This could be due to Snap permissions, if you are using that. **\n\n")
		}

		return nil, err
	}
	defer file.Close()

	ret := make([]string, 0)
	rd := bufio.NewReader(file)

	for {
		bs, p, err := rd.ReadLine()
		if p {
			return nil, fmt.Errorf("line was too long")
		}
		if err == io.EOF {
			return ret, nil
		}
		if err != nil {
			return nil, err
		}
		if len(bs) > 0 {
			ret = append(ret, string(bs))
		}
	}

}
