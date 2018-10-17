// interpol CLI
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/vahidi/interpol"
)

var sep = flag.String("sep", " ", "Column separator")
var lsep = flag.String("lsep", "\n", "Line separator")
var version = flag.Bool("version", false, "Show version information")

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
	// Go flag doesn't go past first "-"
	for _, str := range flag.Args() {
		if str[0] == '-' {
			flag.Usage()
			fail(20, "ERROR: options should be given before commands")
		}
	}

	if *version {
		fmt.Printf("%d.%d.%d\n", interpol.Version[0], interpol.Version[1],
			interpol.Version[2])
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		flag.Usage()
		fail(20, "ERROR: no commands were given")
	}

	// separator strings can contain escaped characters
	sep := unscape(*sep)
	lsep := unscape(*lsep)

	ip := interpol.New()
	strs, err := ip.AddMultiple(flag.Args()...)
	if err != nil {
		fail(20, "ERROR: '%v'", err)
	}

	for {
		for i, s := range strs {
			if i != 0 {
				fmt.Print(sep)
			}
			fmt.Print(s.String())
		}
		fmt.Print(lsep)

		if !ip.Next() {
			break
		}
	}
}
