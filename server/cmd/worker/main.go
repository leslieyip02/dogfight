package main

import (
	"flag"
	"log"
	"server/internal/balancer"
	"server/internal/env"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host := flag.String("host", env.GetOrDefault("HOST", "localhost"), "host")
	port := flag.String("port", env.GetOrDefault("PORT", ":5174"), "port")
	secret := env.GetOrPanic("JWT_SECRET")
	flag.Parse()

	err := balancer.RegisterWorker(*host, *port)
	if err != nil {
		log.Fatalf("could not register worker: %v", err)
	}

	worker := balancer.NewWorker(*host, *port, []byte(secret))
	worker.Serve()
}
