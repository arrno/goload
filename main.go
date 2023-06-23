package main

func main() {
	// hello!
	j := Job{
		Name:        "Job One",
		Workers:     2,
		Interval:    2,
		LogInterval: 5,
		Duration:    12,
		Request: HttpRequest{
			Method:  "GET",
			URI:     "https://www.google.com",
			Body:    nil,
			Headers: map[string]string{},
			Log:     false,
		},
	}
	s := []Job{j}
	runSchedule(s)
}
