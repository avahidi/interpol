package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"bitbucket.org/vahidi/interpol"
)

type item struct {
	By    string
	Title string
	Text  string
}

func downloadPage(url string, agent string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", agent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func showPage(page []byte) error {
	var it item
	err := json.Unmarshal(page, &it)
	if err != nil {
		return err
	}

	if it.Title != "" {
		// self post
		fmt.Printf("\nPOST by %s:\n%s\n%s\n", it.By, it.Title, it.Text)
	} else {
		// comment or link post
		fmt.Printf("\nCOMMENT by %s:\n%s\n", it.By, it.Text)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	ip := interpol.New()
	useragent, err := ip.Add("{{file filename=useragents.txt count=1 mode=random modifier=trim}}")
	if err != nil {
		log.Fatalf("Bad things just happened: %v", err)
	}

	url, err := ip.Add("https://hacker-news.firebaseio.com/v0/item/{{random min=10000 max=999999 count=3}}.json")
	if err != nil {
		log.Fatalf("Other bad things just happened: %v", err)
	}

	fmt.Printf("Today I will be %v\n", useragent)

	for ip.Next() {
		page, err := downloadPage(url.String(), useragent.String())
		if err == nil {
			err = showPage(page)
		}
		if err != nil {
			log.Printf("I guess things can sometimes go wrong: %v", err)
			return
		}
	}
}
