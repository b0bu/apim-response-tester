package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func (c Client) Go(n int) {
	req := c.newRequest()
	httpClient := http.Client{}

	for n > 0 {
		resp, err := httpClient.Do(req)
		log.Printf("job %v created", n)
		if err != nil {
			log.Fatal(err)
		}

		resp.Body.Close()

		go c.Poll(resp)
		n--
	}
}

func (c Client) Poll(r *http.Response) {
	for {
		resp, err := http.Get(r.Header.Get("Operation-Location"))

		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		var p Payload

		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("id: %v status: %v", p.ID, p.Status)
		if p.Status == "complete" {
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

// Progress(id, '.')
func main() {
	numJobs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	c := newClient("https://policy-testing.azure-api.net/api/v1/job/create", http.MethodPost)
	c.Go(numJobs)
}
