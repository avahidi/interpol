// dedicated to Stanislav Petrov who did cloud computing before it was cool
package main

import (
	"bitbucket.org/vahidi/interpol"
	"fmt"
)

// We want to replace 99 with a random number. But since the same number
// is used twice we can't have two random interpolators.
// The solution is to use a random interpolator for the firs occurrence
// and a copy interpolator for the second one to point to the first one
const text = `

Hast du etwas Zeit für mich
Dann singe ich ein Lied für dich
Von {{random min=0 max=99 count=2 format=%2d name=number}} Luftballons
Auf ihrem Weg zum Horizont
Denkst du vielleicht grad an mich
Dann singe ich ein Lied für dich
Von {{copy from=number}} Luftballons
Und, dass so was von so was kommt

`

func main() {
	ip := interpol.New()
	song, err := ip.Add(text)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		for {
			fmt.Printf("%s\n", song.String())
			if !ip.Next() {
				break
			}
		}
	}
}
