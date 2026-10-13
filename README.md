## PulseWatch

PulseWatch, a lightweight uptime monitoring application written in Go.

Monitoring Engine:
Our first milestone is a command-line application that checks multiple websites concurrently and reports their HTTP
status codes and response times.

we separate application entry point from the monitoring logic so that later we can reuse the monitoring package in an
API server or background worker.


### Worker pool.
At this point our app was handling one go routine for every website, which is not ideal as the number of websites grows.
We implement a worker pool to handle this issue, as it uses a fixed number of go routines to process jobs from a shared 
channel


### Automated Testing
Here we are introducing testing 5 scenarious:
    Website returns 200
    Website returns 503
    Website is unreachable
    Website take too long
    Multiple websites.