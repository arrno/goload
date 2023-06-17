package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type METHOD string

const (
	GET    METHOD = "GET"
	POST   METHOD = "POST"
	PUT    METHOD = "PUT"
	PATCH  METHOD = "PATCH"
	DELETE METHOD = "DELETE"
)

type Job struct {
	Name        string
	Workers     uint
	Interval    uint
	LogInterval uint
	Duration    uint
	Request     HttpRequest // should support multiple request types per job
}

type HttpRequest struct {
	Method  METHOD
	URI     string
	Body    any
	Headers map[string]string
	Log     bool
}

func runSchedule(schedule []Job) {
	// allows for ramping up load with series of jobs
	for _, job := range schedule {
		runJob(job)
	}
}

func runJob(j Job) {
	fmt.Println("Running job: " + j.Name)
	var wg sync.WaitGroup
	results := make(chan time.Duration, j.Workers)
	var totalLatency float64
	var seen int
	var spawned int
	for {
		select {
		case <-time.After(time.Duration(j.Duration) * time.Second):
			wg.Wait()
			for seen < spawned {
				re := <-results
				totalLatency += re.Seconds()
				seen++
			}
			fmt.Printf("Total job latency: %v\n", totalLatency)
			fmt.Printf("Average job latency: %v\n", totalLatency/float64(seen))
			close(results)
			return
		case <-time.After(time.Duration(j.Interval) * time.Second):
			wg.Add(1)
			go runRound(j.Workers, j.Request, &wg, results)
			spawned += int(j.Workers)
		case <-time.After(time.Duration(j.LogInterval) * time.Second):
			fmt.Printf("Total job latency: %v\n", totalLatency)
			fmt.Printf("Average job latency: %v\n", totalLatency/float64(seen))
		case re := <-results:
			totalLatency += re.Seconds()
			seen++
		}
	}
}

func runRound(workers uint, req HttpRequest, wg *sync.WaitGroup, results chan<- time.Duration) {
	// spawn x workers to run req
	defer wg.Done()
	var i uint
	for i = 0; i < workers; i++ {
		// run uinit
		go func() {
			r, _ := http.NewRequest(string(req.Method), req.URI, nil)
			startTime := time.Now()
			res, _ := http.DefaultClient.Do(r)
			duration := time.Now().Sub(startTime)
			fmt.Println(req.URI + " " + string(req.Method) + res.Status + duration.String())
			results <- duration
		}()
	}
}
