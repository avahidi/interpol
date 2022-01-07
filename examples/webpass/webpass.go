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

	"github.com/avahidi/interpol"
)

// data for one web access, change this to change the application behavior
type access struct {
	username string
	password string
}

func newAccess(iss []*interpol.InterpolatedString) *access {
	return &access{username: iss[0].String(), password: iss[1].String()}
}

func (s access) values() URL.Values {
	return URL.Values{"username": {s.username}, "password": {s.password}}
}

func (s access) String() string {
	return fmt.Sprintf("'%s'/'%s'", s.username, s.password)
}

// stats
type stats struct {
	start time.Time
	count int
}

func (s *stats) time() string {
	return fmt.Sprintf("%0.2f/s", float32(time.Now().Sub(s.start))/float32(s.count))
}

func (s *stats) report(a *access) {
	s.count++
	if s.count == 1 {
		s.start = time.Now()
	} else if s.count == 10 || s.count == 100 || (s.count%1000) == 0 {
		fmt.Fprintf(os.Stderr, "%8d: %v (%s)...      \r", s.count, a, s.time())
	}
}

func (s *stats) finish() {
	fmt.Printf("\n\nDONE with %d attempts (%s)\n\n", s.count, s.time())
}

// error messages
func failFatal(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("%s (%v)\n", msg, err)
}

func failOnce(err error, count, max int, format string, args ...interface{}) {
	if count >= max {
		failFatal(err, format, args...)
	} else {
		msg := fmt.Sprintf(format, args...)
		log.Printf("%s -- attempt %d/%d failed (%v)...\n", msg, count, max, err)
	}
}

// networking
func post(url string, args URL.Values) ([]byte, error) {
	c := &http.Client{}
	resp, err := c.PostForm(url, args)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Failed with HTTP code %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func tryPost(url string, args URL.Values, tries int) ([]byte, error) {
	err := error(nil)
	for i := 0; i < tries; i++ {
		data, err := post(url, args)
		if err == nil {
			return data, err
		}
		failOnce(err, i+1, tries, url)
		time.Sleep(time.Duration((1 + 3*i)) * time.Second)
	}

	failFatal(err, url)
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

func (s search) report(a *access) {
	fmt.Printf("\n\nFOUND '%v' \n\n", a)
	os.Exit(0)
}

// worker
func worker(id int, s search, wg *sync.WaitGroup, jobs <-chan *access) {
	defer wg.Done()

	for job := range jobs {
		args := job.values()
		resp, err := tryPost(s.url, args, 5)
		if err != nil {
			failFatal(err, "Thread %d failed on %v", id, job)
		}

		text := string(resp)
		if s.verify(text) {
			s.report(job)
			return
		}
	}
}

// parameters
var (
	threads  = flag.Int("threads", 10, "Number of threas")
	url      = flag.String("url", "http://localhost/", "Target URL")
	userdesc = flag.String("username", "", "Username expression, for example {{file filename=usernames.txt}}")
	passdesc = flag.String("password", "", "Password expression")
	postext  = flag.String("postext", "", "Text to check for positive match")
	negtext  = flag.String("negtext", "", "Text to check for positive match")
)

func main() {
	flag.Parse()
	if flag.NArg() != 0 || *userdesc == "" || *passdesc == "" || (*postext == "" && *negtext == "") {
		if *userdesc == "" || *passdesc == "" {
			fmt.Fprintf(os.Stderr, "Missing username or password expressions\n")
		}

		if *postext == "" && *negtext == "" {
			fmt.Fprintf(os.Stderr, "Missing both positive and negative match text\n")
		}
		fmt.Fprintf(os.Stderr, "Usage example:\n%s -url http://localhost/login -username user -password '{{file filename=passwords.txt}}' -postext 'success'\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(3)
	}

	search := search{url: *url, postext: *postext, negtext: *negtext}

	ip := interpol.New()
	vals, err := ip.AddMultiple(*userdesc, *passdesc)
	if err != nil {
		failFatal(err, "Error in interpol expressions")
	}

	var wg sync.WaitGroup
	jobs := make(chan *access, *threads)
	for w := 0; w < *threads; w++ {
		wg.Add(1)
		go worker(w, search, &wg, jobs)
	}

	stat := stats{}
	defer stat.finish()
	for ip.Next() {
		access := newAccess(vals)
		jobs <- access
		stat.report(access)
	}

	close(jobs)
	wg.Wait()
}
