package main

import (
	"log"
	"server/internal/balancer"
)

func main() {
	worker, err := balancer.RegisterWorker("localhost", ":5174")
	if err != nil {
		log.Fatalf("could not register worker: %v", err)
	}
	worker.Serve()
}
