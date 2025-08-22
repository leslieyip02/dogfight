package main

import (
	"log"
	"net/http"
	"os"
	"server/room"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	r.Handle("/*", fs)

	roomManager, err := room.NewManager()
	if err != nil {
		log.Fatalf("failed to create roomManager: %v", err)
		return
	}

	r.Route("/api/room", func(r chi.Router) {
		r.Post("/join", roomManager.HandleJoin)
		r.Get("/players", roomManager.HandleFetchPlayers)
	})
	r.Get("/ws", roomManager.HandleConnect)

	log.Printf("server is running on %s:%s", host, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
