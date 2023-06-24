package main

import (
	"fmt"
	"net/http"
	"strings"
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
	doLog([]string{"Complete...\n"}, true)
}

func runJob(j Job) {

	doLog([]string{"Starting job:", j.Name + "\n"}, true)
	var wg sync.WaitGroup
	results := make(chan time.Duration, j.Workers)
	done := make(chan interface{})
	var totalLatency float64
	var seen int
	var spawned int

	// Max timeout.
	go func() {

		<-time.After(time.Duration(j.Duration) * time.Second)
		wg.Wait()
		for seen < spawned {
			re := <-results
			totalLatency += re.Seconds()
			seen++
		}

		tl := fmt.Sprintf("%v", totalLatency)
		al := fmt.Sprintf("%v", totalLatency/float64(seen))
		doLog([]string{"Total job latency:", truncate(tl, 6), "seconds"}, true)
		doLog([]string{"Average job latency:", truncate(al, 6), "seconds\n"}, false)

		close(results)
		close(done)
	}()
	// Log interval loop.
	go func() {
		for {
			select {
			case <-time.After(time.Duration(j.LogInterval) * time.Second):

				tl := fmt.Sprintf("%v", totalLatency)
				al := fmt.Sprintf("%v", totalLatency/float64(seen))
				doLog([]string{"Total job latency:", truncate(tl, 6), "seconds"}, true)
				doLog([]string{"Average job latency:", truncate(al, 6), "seconds\n"}, false)

			case <-done:
				return
			}
		}
	}()
	// Spawn interval loop.
	for {
		select {
		case <-time.After(time.Duration(j.Interval) * time.Second):
			wg.Add(1)
			go runRound(j.Workers, j.Request, &wg, results)
			spawned += int(j.Workers)

		case re := <-results:
			totalLatency += re.Seconds()
			seen++

		case <-done:
			return
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

			if req.Log {
				d := duration.String()
				doLog([]string{req.URI, string(req.Method), res.Status, truncate(d, 6)}, false)
			}
			results <- duration
		}()
	}
}

// ---- Utilities ---- //

func doLog(args []string, lineBefore bool) {
	var nargs []string
	if lineBefore {
		nargs = append(nargs, "\n"+time.Now().Format("2006-01-02 15:04:05"))
	} else {
		nargs = append(nargs, time.Now().Format("2006-01-02 15:04:05"))
	}
	nargs = append(nargs, args...)
	fmt.Println(strings.Join(nargs, "\t"))
}

func truncate(s string, l int) string {
	if len(s) <= l {
		return s
	}
	for i, _ := range s {
		if i+1 == l {
			return s[:i+1]
		}
	}
	return s
}

func truncateInPlace(s *string, l int) {
	if len(*s) <= l {
		return
	}
	for i, _ := range *s {
		if i+1 == l {
			*s = (*s)[:i+1]
			return
		}
	}
}
