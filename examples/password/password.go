package main

import (
	"bitbucket.org/vahidi/interpol"
	"fmt"
	"log"
)

// assume this functions attempts to login using given credential:
func checkCredential(username string, password string) bool {
	var etcpasswd = map[string]string{
		"paula": "brillant7",
		"alex":  "idreaminexcel",
		"kelly": "secret",
	}
	return etcpasswd[username] == password
}

// Assume this function does whatever sysadms do when they found a weak account
func report(username string) {
	fmt.Printf("Dear %s, please change your password.\n", username)
}

func search(ip *interpol.Interpol, pair []*interpol.InterpolatedString) (found bool) {
	user, password := pair[0], pair[1]
	for {
		if checkCredential(user.String(), password.String()) {
			report(user.String())
			found = true
		}
		if !ip.Next() {
			return
		}
	}

}

func main() {
	// simple check against the files
	ip := interpol.New()
	pair, err := ip.AddMultiple("{{file filename=usernames.txt}}",
		"{{file filename=weakpasswords.txt}}")
	if err != nil {
		log.Fatalf("internal error1: %v", err)
	}

	search(ip, pair)
	// repeat the simple check with additional trailing number
	ip = interpol.New()
	pair, err = ip.AddMultiple("{{file filename=usernames.txt}}",
		"{{file filename=weakpasswords.txt}}{{counter min=0 max=99}}")
	if err != nil {
		log.Fatalf("internal error1: %v", err)
	}

	search(ip, pair)
}
