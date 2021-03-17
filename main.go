package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"fmt"
)


func connect() (redis.Conn, error) {
	return redis.Dial("tcp","localhost:6379")
}

func ListenForJobs(jobs chan Job) {
	// connect to redis
	conn, err := connect()
	if err != nil { panic(err.Error()) }
	defer conn.Close()
	fmt.Println("Waiting for jobs...")

	for {
		reply, err := redis.Values(conn.Do("blpop", "instachat:queue:default", 0))
		if err != nil { panic(err.Error())  }

		var queue string
		var body []byte
		if _, err := redis.Scan(reply, &queue, &body); err != nil {
			fmt.Println(err)
		}
		job := new(Job)

		err = json.Unmarshal(body, &job)
		if err != nil { panic(err.Error()) }

		fmt.Printf("Found job: %s(%s)\n", job.Class, job.Jid)

		jobs <- *job
	}
}

type WorkerFactory func(Job) Worker

func PerformJobs(jobs chan Job) {
	workers := make(map[string]WorkerFactory)

	workers["MessageCreationWorker"] = NewInsetionToDBWorker
	workers["ChatCreationWorker"] = NewInsetionChatToDBWorker

	for {
		// wait for a job

		job := <-jobs

		factory := workers[job.Class]
		worker := factory(job)

		go worker.Perform()
	}
}

func main() {
   for{
	jobs := make(chan Job)
	go ListenForJobs(jobs)
	go PerformJobs(jobs)

	fmt.Println("Press Enter to Exit.")
	var userInput string
	fmt.Scanln(&userInput)
   }
}
