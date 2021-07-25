package main

import (
	"log"
	"time"

	"github.com/joshcarp/swarm"
)

func foo() {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	// Report your test result as a success, if you write it in python, it will looks like this
	// events.request_success.fire(request_type="http", name="foo", response_time=100, response_length=10)
	globalSwarmer.RecordSuccess("http", "foo", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
}

var globalSwarmer *swarm.Swarmer

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	task1 := &swarm.Task{
		Namef:   "foo",
		Weightf: 10,
		Fn:      foo,
	}

	numClients := 10
	spawnRate := 10
	globalSwarmer = swarm.NewStandaloneSwarmer(numClients, float64(spawnRate))
	globalSwarmer.Run(task1)
}
