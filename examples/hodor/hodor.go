// this example demonstrates creation of custom handlers
//
// disclaimer: I have not seen a single episode of GoT. I live in Sweden,
//             and don't need to be reminded that winter is coming.

package main

import (
	"fmt"
	"log"

	"github.com/avahidi/interpol"
)

// this is the internal structure the handler will use for its own data
type hodor struct {
	hodor int
	count int
}

// this is a handler factory function, called to create new instances
func hodorHodor(ctx *interpol.Interpol, text string, data *interpol.Parameters) (interpol.Handler, error) {
	ret := &hodor{
		count: data.GetInteger("count", 5),
	}
	ret.Reset()
	return ret, nil
}

// these implements the interpol.Handler interface
func (hodor *hodor) String() string {
	if hodor.hodor > 1 {
		return "hodor, "
	}
	return "hodor!"
}

func (hodor *hodor) Next() bool {
	// note how we check twice
	if hodor.hodor > 0 {
		hodor.hodor--
	}
	return hodor.hodor > 0
}

func (hodor *hodor) Reset() {
	hodor.hodor = hodor.count
}

// the main function of this examples shows how to register a new handler
func main() {

	ip := interpol.New()
	ip.AddHandler("hodor", hodorHodor)

	hodor, err := ip.Add("{{hodor count=4}}")
	if err != nil {
		log.Fatalf("Failed: %v\n", err)
	}

	for ip.Next() {
		fmt.Printf(hodor.String())
	}
	fmt.Println()

}
