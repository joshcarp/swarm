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

func waitForQuit() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	quitByMe := false
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		quitByMe = true
		globalBoomer.Quit()
		wg.Done()
	}()

	swarm.Events.Subscribe("boomer:quit", func() {
		if !quitByMe {
			wg.Done()
		}
	})

	wg.Wait()
}

var globalBoomer = swarm.NewBoomer("127.0.0.1", 5557)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ts := swarm.NewWeighingTaskSet()

	taskA := &swarm.Task{
		Namef:   "TaskA",
		Weightf: 10,
		Fn: func() {
			time.Sleep(100 * time.Millisecond)
			globalBoomer.RecordSuccess("task", "A", 100, int64(10))
		},
	}

	taskB := &swarm.Task{
		Namef:   "TaskB",
		Weightf: 20,
		Fn: func() {
			time.Sleep(100 * time.Millisecond)
			globalBoomer.RecordSuccess("task", "B", 100, int64(20))
		},
	}

	// Expecting RPS(taskA)/RPS(taskB) to be close to 10/20
	ts.AddTask(taskA)
	ts.AddTask(taskB)

	task := &swarm.Task{
		Namef: "TaskSet",
		Fn:    ts.Run,
	}

	globalBoomer.Run(task)

	waitForQuit()
	log.Println("shut down")
}
