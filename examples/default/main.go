package main

import (
	"log"
	"time"

	"github.com/joshcarp/swarm"
)

func foo(bm *swarm.Swarmer) func() {
	return func() {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		elapsed := time.Since(start)

		// Report your test result as a success, if you write it in python, it will looks like this
		// events.request_success.fire(request_type="http", name="foo", response_time=100, response_length=10)
		bm.RecordSuccess("http", "foo", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
	}
}

func bar(bm *swarm.Swarmer) func() {
	return func() {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		elapsed := time.Since(start)

		// Report your test result as a failure, if you write it in python, it will looks like this
		// events.request_failure.fire(request_type="udp", name="bar", response_time=100, exception=Exception("udp error"))
		bm.RecordFailure("udp", "bar", elapsed.Nanoseconds()/int64(time.Millisecond), "udp error")
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	bm := swarm.NewSwarmer("localhost", 5557)

	task1 := &swarm.Task{
		Namef:   "foo",
		Weightf: 10,
		Fn:      foo(bm),
	}

	task2 := &swarm.Task{
		Namef:   "bar",
		Weightf: 30,
		Fn:      bar(bm),
	}
	bm.Run(task1, task2)
}
