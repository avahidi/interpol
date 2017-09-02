// pocli = run interpol from cli
package main

import (
	"bitbucket.org/vahidi/interpol"
	"flag"
	"fmt"
	"log"
	"strings"
)

var defaultInput = []string{
	"{{set count=4 sep=',' data='Every single day,Every word you say,Every game you play,Every night you stay'}}",
	"I'll be watching you",
}

var separator = flag.String("sep", "\n", "separator")

func unscape(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\t", "\t", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	return s
}

func parseParameters() []string {
	flag.Parse()
	if flag.NArg() == 0 {
		return defaultInput
	}
	return flag.Args()
}

func main() {
	cmds := parseParameters()
	sep := unscape(*separator)

	ip := interpol.New()
	strs, err := ip.AddMultiple(cmds...)
	if err != nil {
		log.Fatalf("Officer down: %v\n", err)
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
