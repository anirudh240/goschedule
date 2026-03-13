# GoSchedule

## What is this?
(one paragraph — what does this project do, explain it like you're telling a friend)
A go scheduler used to prioritize jobs or run jobs concurrently or parallely 

## Why I built this
I built this to get a better understanding of GOLANG, this is my first ever project, 
this project helps me understand how API's work and also the scheduling bit. 
The concept of error handling and the difference between concurrency and parallelism 
I was actively contributing to a cncf project named VOLCANO, it is a Kubernetes batch scheduling system
it was built to extend and enhance the capabilities kube-scheduler


## Features
- Priority based job scheduling
- Worker pool with concurrency control
- HTTP API to view and add jobs
- Error handling with fail fast principle
- JSON based job submission
- Containerized with Docker

## How to run it

### Normal:
go run main.go

### With Docker:
docker build -t scheduler .
docker run -p 9090:9090 scheduler

## API
GET /jobs - returns all jobs and their status
POST /add - add a new job
sends a JSON body like: 
{
    "name": "job name",
    "priority": 1,
    "error": false
}
