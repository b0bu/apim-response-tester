package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	URL     string
	Headers map[string]string
	Method  string
}

type Payload struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c Client) newRequest() *http.Request {
	req, err := http.NewRequest(c.Method, c.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	return req
}

func (c Client) Go() {
	req := c.newRequest()
	httpClient := http.Client{}
	var wg sync.WaitGroup

	Clear()

	for i := 1; i <= state.MaxLines; {
		wg.Add(1)
		resp, err := httpClient.Do(req)
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
		resp, err := http.Get(r.Header.Get("Operation-Location"))
		if err != nil {
			log.Fatalln(err)
		}

		var job Payload
		err = json.NewDecoder(resp.Body).Decode(&job)
		if err != nil {
			log.Fatalln(err)
		}

		if err != nil {
			log.Fatalln(err)
		}
		Progress(threadID, fmt.Sprintf("job id: %v status: %v\n", job.ID, job.Status))
		if job.Status == "complete" {
			Progress(threadID, fmt.Sprintf("job id: %v status: %v\n", job.ID, job.Status))
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

func newClient(url, method string) Client {
	return Client{
		url,
		map[string]string{"Ocp-Apim-Subscription-Key": getToken(), "Content-Length": "0"},
		method,
	}
}

func main() {
	drawlines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	state.MaxLines = drawlines

	c := newClient("https://policy-testing.azure-api.net/api/v1/job/create", http.MethodPost)
	c.Go()
}
