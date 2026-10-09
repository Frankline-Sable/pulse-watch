## PulseWatch

PulseWatch, a lightweight uptime monitoring application written in Go.

Monitoring Engine:
Our first milestone is a command-line application that checks multiple websites concurrently and reports their HTTP
status codes and response times.

we separate application entry point from the monitoring logic so that later we can reuse the monitoring package in an
API server or background worker.