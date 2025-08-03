package main

import (
	"log"
	"net/http"
	"os"
)

func pollJobRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	tokenFile := ".token"
	bytes, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = http.Header{
		"Ocp-Apim-Subscription-Key": {string(bytes)},
	}
	return req
}

/*
num jobs 1 by default
create the jobs
watch them using the operation-location url which should be rewritten via apim

*/

func main() {

	c := http.Client{}

	numJobs := os.Args[1]

	url := "https://policy-testing.azure-api.net/api/v1/job/" + id

	req := pollJobRequest(url)

	for {
		resp, err := c.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		// check status
	}

}
