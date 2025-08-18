package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	client  *http.Client
	Headers map[string]string
}

type Payload struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Backend string `json:"backend"`
}

func (c Client) newRequest(url string, method string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	return req
}

func (c Client) Go() {
	req := c.newRequest("https://policy-testing.azure-api.net/api/v1/job/create", http.MethodPost)
	var wg sync.WaitGroup

	Clear()

	for i := 1; i <= state.MaxLines; {
		wg.Add(1)
		resp, err := c.client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		go func(id int, r *http.Response) {
			c.Poll(id, r)
			wg.Done()
			defer resp.Body.Close()
		}(i, resp)
		i++
	}

	wg.Wait()
	defer Return()
}

func (c Client) Poll(threadID int, r *http.Response) {
	for {
		req := c.newRequest(r.Header.Get("Operation-Location"), http.MethodGet)
		resp, err := c.client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		var message Payload
		err = json.NewDecoder(resp.Body).Decode(&message)
		if err != nil {
			log.Fatalln(err)
		}

		Progress(threadID, message)
		if message.Status == "complete" {
			Progress(threadID, message)
			return
		}
	}
}

func getToken() string {
	tokenFile := ".token"
	bytes, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.TrimRight(string(bytes), "\n")
}

func newClient() Client {
	return Client{
		&http.Client{
			Timeout: time.Second * 60,
		},
		map[string]string{"Ocp-Apim-Subscription-Key": getToken(), "Content-Length": "0"},
	}
}

func main() {
	drawlines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	state.MaxLines = drawlines

	c := newClient()
	c.Go()
}
