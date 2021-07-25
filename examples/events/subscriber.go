package main

import (
	"log"
	"time"

	"github.com/joshcarp/swarm"
)

// This is an example about how to subscribe to boomer's internal events.

func foo(bm *swarm.Boomer) func(){
	return func(){
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		elapsed := time.Since(start)
		bm.RecordSuccess("http", "foo", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	bm := swarm.NewBoomer("localhost", 5557)
	bm.Events.Subscribe("boomer:spawn", func(workers int, spawnRate float64) {
		log.Println("The master asks me to spawn", workers, "goroutines with a spawn rate of", spawnRate, "per second.")
	})

	bm.Events.Subscribe("boomer:stop", func() {
		log.Println("The master asks me to stop.")
	})

	bm.Events.Subscribe("boomer:quit", func() {
		log.Println("Boomer is quitting now, may be the master asks it to do so, or it receives one of SIGINT and SIGTERM.")
	})

	task := &swarm.Task{
		Namef:   "foo",
		Weightf: 10,
		Fn:      foo(bm),
	}

	bm.Run(task)
}
