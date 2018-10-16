// interpol CLI
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/vahidi/interpol"
)

var sep = flag.String("sep", " ", "Separator")

func unscape(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\t", "\t", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	return s
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <commands>  \n", os.Args[0])
		flag.PrintDefaults()
	}
}
func fail(code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(code)
}
func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		fail(20, "ERROR: no commands were given")

	}

	sep := unscape(*sep)

	ip := interpol.New()
	strs, err := ip.AddMultiple(flag.Args()...)
	if err != nil {
		fail(20, "ERROR: '%v'", err)
	}

	for {
		for _, s := range strs {
			fmt.Printf("%s%s", s.String(), sep)
		}
		fmt.Println()

		if !ip.Next() {
			break
		}
	}
}
