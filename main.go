package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Job struct {
	ID string `json:"id"`
}

var jobs = []Job{}

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

func createJob(c *gin.Context) {
	var job Job
	uid := strconv.Itoa(uuid())
	jobs = append(jobs, Job{ID: uid})
	c.Header("operation-location", "http://10.10.10.10/my/long/path/job/"+uid)
	c.JSON(http.StatusCreated, job)
}

func main() {
	router := gin.Default()
	router.POST("/job/create", createJob)
	router.GET("/job/:id", getJob)
	router.GET("/jobs", getJobs)

	router.Run("localhost:8080")
}
