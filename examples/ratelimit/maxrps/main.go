package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func waitForQuit() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		globalSwarmer.Quit()
		wg.Done()
	}()

	wg.Wait()
}

var globalSwarmer = swarm.NewSwarmer("127.0.0.1", 5557)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	task1 := &swarm.Task{
		Namef:   "foo",
		Weightf: 10,
		Fn:      foo,
	}

	ratelimiter := swarm.NewStableRateLimiter(100, time.Second)
	log.Println("the max rps is limited to 100/s.")
	globalSwarmer.SetRateLimiter(ratelimiter)

	globalSwarmer.Run(task1)

	waitForQuit()
	log.Println("shut down")
}
