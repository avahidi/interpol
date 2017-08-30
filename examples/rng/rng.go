package main

import (
	"bitbucket.org/vahidi/interpol"
	"fmt"
	"log"
)

func main() {
	ip := interpol.New()
	rs, err := ip.Add("{{random min=0 max=9 count=3}}{{random min=0 max=999 count=3 format=%03d}}")
	if err != nil {
		log.Fatalf("Bad things just happened: %v", err)
	}

	fmt.Printf("My friend, quality random, only for you: ")
	for {
		fmt.Printf("%s ", rs)
		if !ip.Next() {
			break
		}
	}
	fmt.Println()
}
