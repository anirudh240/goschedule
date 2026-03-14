# GoSchedule

A Go based job scheduler that prioritizes and runs jobs concurrently using a worker pool.

## Why I Built This

While actively contributing to Volcano — a CNCF batch scheduling system that extends and enhances 
Kubernetes native batch Scheduling  — I wanted to understand and explore the core concepts behind 
job scheduling. This project is my attempt to build a simplified version 
of what Volcano does under the hood. It helped me understand how APIs work, how 
concurrency differs from parallelism, and how error handling works in real systems.

## What I Learned
- How goroutines and channels work in Go
- The difference between concurrency and parallelism
- How Docker containerization works
- How HTTP APIs are structured and how JSON bridges different systems

## Features
- Priority based job scheduling
- Worker pool with concurrency control
- HTTP API to view and add jobs
- Error handling with fail fast principle
- JSON based job submission
- Containerized with Docker

## How to Run

**Normal:**
```bash
go run main.go
```

**With Docker:**
```bash
docker build -t scheduler .
docker run -p 9090:9090 scheduler
```

## API

**GET /jobs** - returns all jobs and their current status

**POST /add** - add a new job by sending a JSON body:
```json
{
    "name": "job name",
    "priority": 1,
    "error": false
}
```
Lower priority number runs first.


## Kubernetes Deployment

This project can be deployed on a Kubernetes cluster using the provided manifests.

### Prerequisites
- Docker
- Minikube
- kubectl

### Deploy locally with Minikube

Start your cluster:
```bash
minikube start
```

Load the Docker image into minikube:
```bash
minikube image load scheduler
```

Apply the manifests:
```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

Access the scheduler:
```bash
minikube service goschedule-service
```

Then hit `/jobs` on the URL minikube gives you.

### What Kubernetes does here
- Runs 2 replicas of the scheduler simultaneously
- Automatically restarts pods if they crash
- Exposes the scheduler via a Service
- Maintains desired state forever — delete a pod and Kubernetes spins up a new one instantly
