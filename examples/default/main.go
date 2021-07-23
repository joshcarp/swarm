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
	swarm.RecordSuccess("http", "foo", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
}

func bar() {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	// Report your test result as a failure, if you write it in python, it will looks like this
	// events.request_failure.fire(request_type="udp", name="bar", response_time=100, exception=Exception("udp error"))
	swarm.RecordFailure("udp", "bar", elapsed.Nanoseconds()/int64(time.Millisecond), "udp error")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	task1 := &swarm.Task{
		Namef:   "foo",
		Weightf: 10,
		Fn:      foo,
	}

	task2 := &swarm.Task{
		Namef:   "bar",
		Weightf: 30,
		Fn:      bar,
	}

	swarm.Run(task1, task2)
}
