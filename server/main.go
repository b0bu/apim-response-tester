package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Job struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Backend string `json:"backend"`
}

var jobs = []*Job{}

/*
create jobs
get jobs
get job

set header with bad ip and path
use apim c# policy to rewrite header to client
use apim c# policy to forward get operation back
*/

func uuid() int {
	return rand.Int()
}

func getJobs(c *gin.Context) {
	c.JSON(http.StatusOK, jobs)
}

func getJob(c *gin.Context) {
	id := c.Param("id")

	for _, job := range jobs {
		if job.ID == id {
			c.JSON(http.StatusOK, job)
			return
		}
	}
}

func work(j *Job) {
	n := rand.Intn(30)
	time.Sleep(time.Duration(n) * time.Second)
	j.Status = "complete"
}

func createJob(c *gin.Context) {

	backend := os.Getenv("BACKEND_ID")

	uid := strconv.Itoa(uuid())
	job := &Job{ID: uid, Status: "pending", Backend: backend}
	jobs = append(jobs, job)

	endpoint := os.Getenv("ENDPOINT_FQDN")

	if endpoint == "" {
		endpoint = "localhost"
	}

	go work(job)
	c.Header("operation-location", "http://"+endpoint+":8080/job/"+uid)
	c.JSON(http.StatusCreated, *job)
}

func main() {
	router := gin.Default()
	router.POST("/job/create", createJob)
	router.GET("/job/:id", getJob)
	router.GET("/jobs", getJobs)

	router.Run("0.0.0.0:8080")
}
