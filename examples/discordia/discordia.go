// This example demonstrates creation of custom modifiers. We will create
// the "discordia" which will inject some "text" about "rate" percent
// of times.
//

package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/avahidi/interpol"
)

// input data stolen from wikipedia:
const input = "Flower power was a slogan used during the late 1960s and early 1970s as a symbol of passive resistance and non-violence ideology."

// this is the internal structure for the modifier
type discordiaModifier struct {
	rate int
	text string
}

func (m *discordiaModifier) Modify(str string) string {
	if (rand.Int() % 100) < m.rate {
		return fmt.Sprintf("%s %s", str, m.text)
	}
	return str
}

// interpol will use this function to create a modifier.
// Note how we extract parameters from the input
func newDiscordiaModifier(ctx *interpol.Interpol, data *interpol.Parameters) (interpol.Modifier, error) {
	return &discordiaModifier{
		rate: data.GetInteger("modifier-rate", 10),
		text: data.GetString("modifier-text", "emc2"),
	}, nil
}

// the main function of this examples shows how to register a new modifier
func main() {

	// create a context and register our discordia modifier with it
	ip := interpol.New()
	ip.AddModifier("discordia", newDiscordiaModifier)

	// with the discordia in place, we can use it in an interpolation
	cmd := fmt.Sprintf("{{set sep=' ' data='%s' modifier=discordia modifier-text=fnord modifier-rate=23}}", input)
	str, err := ip.Add(cmd)
	if err != nil {
		log.Fatalf("Something bad happened: %v", err)
	}

	for ip.Next() {
		fmt.Printf("%s ", str.String())
	}
	fmt.Println()
}
