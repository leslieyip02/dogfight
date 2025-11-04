package main

import (
	"server/internal/balancer"
)

func main() {
	master := balancer.NewMaster("http://localhost", ":5173")
	master.Serve()
}
