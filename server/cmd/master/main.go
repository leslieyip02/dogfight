package main

import (
	"flag"
	"server/internal/balancer"
	"server/internal/env"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host := flag.String("host", env.GetOrDefault("HOST", "localhost"), "host")
	port := flag.String("port", env.GetOrDefault("PORT", ":5173"), "port")
	secret := env.GetOrPanic("JWT_SECRET")
	flag.Parse()

	master := balancer.NewMaster(*host, *port, []byte(secret))
	master.Serve()
}
