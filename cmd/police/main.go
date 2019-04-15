// interpol CLI
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"bitbucket.org/vahidi/interpol"
)

var file = flag.String("f", "", "Read commands from this file")
var sep = flag.String("sep", " ", "Column separator")
var lsep = flag.String("lsep", "\n", "Line separator")
var version = flag.Bool("version", false, "Show version information")
var seed = flag.Int64("seed", 0, "Random number generator seed (0 means use system time)")

func unscape(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\t", "\t", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	return s
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Police is a command-line interface for Interpol, "+
			"which is a string interpolation used for penetration testing, fuzzing, and much more."+
			"\n\n"+
			"Usage: %s [options] <commands>  \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n"+
			"\tpolice -sep \", \" \"Hello\" \"{{set sep=' ' data='Kitty World Dolly goodbye'}}!\"\n"+
			"\tpolice -lsep \":\" \"{{random min=0 max=255 count=8 format=%%02x}}\"\n")
	}
}
func fail(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\n"+format+"\n", a...)
	os.Exit(code)
}

func main() {
	flag.Parse()
	commands := flag.Args()

	// Go flag doesn't go past first "-"
	for _, str := range commands {
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

	// see if there are any commands we should read from file first
	if *file != "" {
		commandsFromFile, err := interpol.ReadFile(*file)
		if err != nil {
			fail(20, "ERROR: could not read commands from file - %v", err)
		}
		commands = append(commands, commandsFromFile...)
	}

	if len(commands) == 0 {
		flag.Usage()
		fail(20, "ERROR: no commands were given")
	}

	// set random seed
	if *seed != 0 {
		rand.Seed(*seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	// separator strings can contain escaped characters
	sep := unscape(*sep)
	lsep := unscape(*lsep)

	ip := interpol.New()
	strs, err := ip.AddMultiple(commands...)
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
