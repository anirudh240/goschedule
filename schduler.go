package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Job struct {
	Name     string
	ID       int
	Status   string
	Priority int
	Error    bool // Error handling
}

// the function below is for the channels

func worker(jobChannel chan *Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobChannel {

		// error handling situation

		if job.Error {
			job.Status = "Failed"
			fmt.Printf("Error: Job #%d %s could not complete\n", job.ID, job.Name)
			continue
		}

		fmt.Printf("Working on job #%d: %s | Priority: %d\n", job.ID, job.Name, job.Priority)
		time.Sleep(2 * time.Second)

		job.Status = "Completed"
		fmt.Printf("Success: %s is now %s\n", job.Name, job.Status)
	}
}

func main() {
	queue := []Job{
		{Name: "NVMe Setup", ID: 2, Status: "Pending", Priority: 2, Error: true},
		{Name: "Calibrate update", ID: 3, Status: "Pending", Priority: 3},
		{Name: "POP_OS_UPDATE", ID: 1, Status: "Pending", Priority: 1},
	}
	fmt.Println("Firing up the scheduler......")
	time.Sleep(2 * time.Second)

	var wg sync.WaitGroup
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Priority < queue[j].Priority
	})

	// HTTP method which exposes this program to API

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // w is what we write to the server and r is what we recieve
		fmt.Fprintln(w, "Scheduler is running")
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		for _, job := range queue {
			fmt.Fprintf(w, "Job #%d: %s | Priority: %d | Status: %s\n",
				job.ID, job.Name, job.Priority, job.Status)
		}
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var job Job
		err := json.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		queue = append(queue, job)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Job added: %s", job.Name)
	})

	go http.ListenAndServe(":9090", nil)

	// adding of channels here

	jobChannel := make(chan *Job, len(queue))

	// feeds the jobs into channels

	for i := 0; i < 2; i++ { // set the limit to 1 and its sequential and if we put it to 3 we get concurrent
		wg.Add(1)
		go worker(jobChannel, &wg)
	}

	for i := range queue {
		jobChannel <- &queue[i]
	}
	close(jobChannel)

	wg.Wait()
	fmt.Println("All the given tasks are finished")

	select {}
}

// TO GET the best understanding of the overall code, try adding wg.wait() and the print statement insdie the for loop
// TRY ADDING wg.ADD inside the first function

// adding priorities in this queue
// we would start with making changes inside the struct

//then we add priority inside the func main where we did
// queue := []Job {}

// this still ends up running concurrently(jobs running without prio)

//concurrency is used for speed and concurrency is not parallelism
// sequential prio is used for hospital systems etc

//better and the last example is a web server which tells "handle all requests concurrently but paid users get priority"

//for _, job := range queue copies the data from slice into job variable and so when the status changes to completed it never changes because its copied

// for i := range changes the original thing which is better in this case cuz we are changing the status manually later

// Q1: What happens if you put wg.Add(1) inside processJob instead of main — why would that break things?
// Q2: You wrote &queue[i] — what does the & do and what would happen if you just wrote queue[i]?
// Q3: Your sort uses queue[i].Priority < queue[j].Priority — if you changed < to > what changes in the output?
// Q4: Why does wg.Wait() being inside the loop make it sequential and outside make it concurrent — explain it in your own words like you're explaining to a friend?
//
//
// wg.Add(1) if we put it inside the process job, the program would run but the program doesnt know waht to increment like the output would be 'all tasks finsihed" and one rroutine would run but nothing else would get printed and the program woud likely end (RACE CONDITION)
//
// when we are writing & here we are making sure go doesnt make a copy of the queue and make changes there, when writing & we are making sure everything comes back to the same address
//
// i think nothing would change with < and > considering the outputs were always random, this is a concurrency inclined program and not a sequential prio program
//
// wg.Wait() if we add inside the for loop, the loop waits after printing evry task once and then prints the end too but it waits after every iteration which is not the desired output
