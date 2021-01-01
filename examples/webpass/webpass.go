package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	URL "net/url"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/vahidi/interpol"
)

// parameters
var threads = flag.Int("threads", 10, "Number of threas")
var url = flag.String("url", "http://localhost/", "Target URL")
var userdesc = flag.String("username", "", "Username expression, for example {{file filename=usernames.txt}}")
var passdesc = flag.String("password", "", "Password expression")
var postext = flag.String("postext", "", "Text to check for positive match")
var negtext = flag.String("negtext", "", "Text to check for positive match")

// error messages
func failOnce(url string, count, max int, err error) {
	fmt.Fprintf(os.Stderr, "failure %d/%d (%v)...\n", count, max, err)
}

func failFatal(url string, err error) {
	fmt.Fprintf(os.Stderr, "Final failure %v\n", err)
	os.Exit(20)
}

// networking
func post(url string, args URL.Values) ([]byte, error) {
	c := &http.Client{}

	resp, err := c.PostForm(url, args)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func tryPost(url string, args URL.Values, tries int) ([]byte, error) {
	err := error(nil)
	for i := 0; i < tries; i++ {
		data, err := post(url, args)
		if err == nil {
			return data, err
		}
		failOnce(url, i+1, tries, err)
		time.Sleep(time.Duration((1 + 3*i)) * time.Second)
	}

	failFatal(url, err)
	return nil, err
}

// search
type search struct {
	url     string
	postext string
	negtext string
}

func (s search) verify(text string) bool {
	return (s.negtext != "" && !strings.Contains(string(text), s.negtext)) ||
		(s.postext != "" && strings.Contains(string(text), s.postext))
}

func (s search) report(username, password string) {
	fmt.Printf("FOUND '%s'/'%s'\n", username, password)
	os.Exit(0)
}

// worker
type job struct {
	username string
	password string
}

func worker(s search, wg sync.WaitGroup, jobs <-chan job) {
	defer wg.Done()

	for job := range jobs {
		args := URL.Values{"username": {job.username}, "password": {job.password}}
		resp, err := tryPost(s.url, args, 5)
		if err != nil {
			fmt.Printf("FAILED on %v: %v\n", job, err)
			os.Exit(3)
		}

		text := string(resp)
		if s.verify(text) {
			s.report(job.username, job.password)
			return
		}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 0 || *userdesc == "" || *passdesc == "" || (*postext == "" && *negtext == "") {
		fmt.Fprintf(os.Stderr, "Usage: %s -url <url> -username <expression> -password <expression> -postext/-negtext <text>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(3)
	}

	s := search{url: *url, postext: *postext, negtext: *negtext}

	ip := interpol.New()
	vals, err := ip.AddMultiple(*userdesc, *passdesc)
	if err != nil {
		log.Fatalf("Bad things just happened: %v", err)
	}

	var wg sync.WaitGroup
	jobs := make(chan job, *threads)
	for w := 1; w < *threads; w++ {
		wg.Add(1)
		go worker(s, wg, jobs)
	}

	i := 0
	for ip.Next() {
		username := vals[0].String()
		password := vals[1].String()
		jobs <- job{username: username, password: password}

		i = i + 1
		if (i % 1000) == 0 {
			fmt.Printf("%d: %s %s...          \r", i, username, password)
		}
	}

	close(jobs)
	wg.Wait()
}
