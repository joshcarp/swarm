package main

import (
	"context"
	"log"
	"time"

	"github.com/joshcarp/swarm"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var globalClient *clientv3.Client

func worker(bm *swarm.Swarmer) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		start := time.Now()
		resp, err := globalClient.Put(ctx, "hello", "swarmer")
		elapsed := time.Since(start)
		if err != nil {
			bm.RecordFailure("etcd", "put", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		} else {
			bm.RecordSuccess("etcd", "put", elapsed.Nanoseconds()/int64(time.Millisecond), int64(resp.Header.Size()))
		}

		cancel()
	}
}

func main() {
	client, err := clientv3.NewFromURL("127.0.0.1:2379")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	globalClient = client
	bm := swarm.NewSwarmer("localhost", 5557)
	task := &swarm.Task{
		Namef: "etcd/clientv3",
		Fn:    worker(bm),
	}

	bm.Run(task)
}
