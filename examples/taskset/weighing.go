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

func waitForQuit(bm *swarm.Swarmer) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	quitByMe := false
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		quitByMe = true
		bm.Quit()
		wg.Done()
	}()

	bm.Bus.Subscribe(swarm.EventQuit, func() {
		if !quitByMe {
			wg.Done()
		}
	})

	wg.Wait()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	bm := swarm.NewSwarmer("127.0.0.1", 5557)
	ts := swarm.NewWeighingTaskSet()

	taskA := &swarm.Task{
		Namef:   "TaskA",
		Weightf: 10,
		Fn: func() {
			time.Sleep(100 * time.Millisecond)
			bm.RecordSuccess("task", "A", 100, int64(10))
		},
	}

	taskB := &swarm.Task{
		Namef:   "TaskB",
		Weightf: 20,
		Fn: func() {
			time.Sleep(100 * time.Millisecond)
			bm.RecordSuccess("task", "B", 100, int64(20))
		},
	}

	// Expecting RPS(taskA)/RPS(taskB) to be close to 10/20
	ts.AddTask(taskA)
	ts.AddTask(taskB)

	task := &swarm.Task{
		Namef: "TaskSet",
		Fn:    ts.Run,
	}

	bm.Run(task)

	waitForQuit(bm)
	log.Println("shut down")
}
