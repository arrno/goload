# goload
A humble http stress tester written in Go to make use of the concurrency model.
## Example
Initiate a job
```Go
j := Job{
    Name:        "Load Test One",
    Workers:     2,
    Interval:    2,
    LogInterval: 5,
    Duration:    12,
    Request: HttpRequest{
        Method:  "GET",
        URI:     "https://www.example.com",
        Body:    nil,
        Headers: map[string]string{},
        Log:     true,
    },
}
s := []Job{j}
runSchedule(s)
```

Evaluate results
```
2023-06-24 17:41:21     Starting job:   Load Test One

2023-06-24 17:41:24     https://www.example.com  GET     200 OK  408.56
2023-06-24 17:41:24     https://www.example.com  GET     200 OK  410.16
2023-06-24 17:41:26     https://www.example.com  GET     200 OK  352.12
2023-06-24 17:41:26     https://www.example.com  GET     200 OK  353.91

2023-06-24 17:41:26     Total job latency:      1.5247  seconds
2023-06-24 17:41:26     Average job latency:    0.3811  seconds

2023-06-24 17:41:28     https://www.example.com  GET     200 OK  250.11
2023-06-24 17:41:28     https://www.example.com  GET     200 OK  257.98
2023-06-24 17:41:31     https://www.example.com  GET     200 OK  347.55
2023-06-24 17:41:31     https://www.example.com  GET     200 OK  349.34

2023-06-24 17:41:31     Total job latency:      2.7297  seconds
2023-06-24 17:41:31     Average job latency:    0.3412  seconds

2023-06-24 17:41:33     https://www.example.com  GET     200 OK  258.72
2023-06-24 17:41:33     https://www.example.com  GET     200 OK  259.97

2023-06-24 17:41:33     Total job latency:      3.2484  seconds
2023-06-24 17:41:33     Average job latency:    0.3248  seconds

2023-06-24 17:41:34     Complete...
```