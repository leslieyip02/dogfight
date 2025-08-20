package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host, found := os.LookupEnv("HOST")
	if !found {
		host = "http://localhost"
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		port = "3000"
	}

	buildDir := http.Dir("../client/dist")
	fs := http.FileServer(buildDir)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/*", fs)

	// TODO: handle websockets

	log.Printf("server is running on %s:%s", host, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
