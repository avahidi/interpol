Interpol can be used as a library in your own programs. The following snippet highlights the core parts of the library::


    package main

    import (
        "fmt"
        "log"
        "github.com/avahidi/interpol"
    )

    func main() {
        ip := interpol.New()
        vs, err := ip.AddMultiple(
            "{{counter min=10 max=33 step=7}}",
            "{{set data='ABCD' mode='linear'}}",
            "{{counter min=0 max=9}}",
        )

        if err != nil {
            log.Fatal(err)
        }

        for ip.Next() {
            fmt.Printf("%s-%s-%s\n", vs[0], vs[1], vs[2])
        }
    }

This code will generate the following output::

    10-A-0
    17-A-0
    ...
    31-D-)


More examples
~~~~~~~~~~~~~

For more advanced examples, see the following folders:

- **hodor** - demonstrates use of custom interpolators (as the name clearly implies :) )
- **discordia** - demonstrates use of custom modifiers
- **nena** - demonstrates use of copy
- **hackernews** - download 3 random HN comments from firebase using a random user-agent
- **webpass** - web form brute force example, because we are too cool to use hydra

