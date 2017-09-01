// Interpol is a string interpolation library for generating
// lists of strings out of a set of rules
//
// A interpolation series is defined by a contex and a number of strings:
//
//    ctx := interpol.New()
//    str1, err := ctx.Add("{{file filename=/etc/passwd}}")
//    str2, err := ctx.Add("{{set data=good|bad sep=| mode=random count=1}}")
//
// To iterate over all values use Next(), to read a value use String()
//
//    for {
//        fmt.Printf("str1=%s str2=%s\n", str1.String(), str2.String()
//        if ! ctx.Next() {
//            break
//        }
//    }

package interpol
