package balancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/internal/room"
	"server/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

var corsHandler = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"http://localhost:5173"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	AllowCredentials: true,
})

type Worker struct {
	host string
	port string

	lobby  *room.Lobby
	secret []byte
}

func NewWorker(host string, port string) *Worker {
	godotenv.Load()
	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		log.Fatalf("JWT_SECRET must be set")
	}

	return &Worker{
		host:   host,
		port:   port,
		lobby:  room.NewLobby(),
		secret: []byte(secret),
	}
}

func RegisterWorker(host string, port string) (*Worker, error) {
	body, err := json.Marshal(RegisterRequest{
		Host: host,
		Port: port,
	})
	if err != nil {
		return nil, err
	}

	url := "http://localhost:5173/internal/register"
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %v", response.StatusCode)
	}

	return NewWorker(host, port), nil
}

func (w *Worker) Serve() {
	r := chi.NewRouter()
	r.Use(corsHandler)
	r.Use(middleware.Logger)

	r.Get("/api/room/snapshot", w.HandleSnapshot)
	r.Get("/api/room/ws", w.HandleWS)
	r.Put("/internal/create", w.HandleCreate)
	r.Get("/internal/status", w.HandleStatus)

	log.Printf("game server is running on http://%s%s", w.host, w.port)
	if err := http.ListenAndServe(w.port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (w *Worker) HandleSnapshot(rw http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := session.ParseToken(token, w.secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	roomId := claims["roomId"].(string)
	snapshot := w.lobby.GetSnapshot(roomId)
	if snapshot == nil {
		http.Error(rw, fmt.Sprintf("could not find room %s", roomId), http.StatusNotFound)
		return
	}

	body, err := proto.Marshal(snapshot)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	if _, err = rw.Write(body); err != nil {
		message := "unable to get room state"
		http.Error(rw, message, http.StatusInternalServerError)
		return
	}
}

func (w *Worker) HandleWS(rw http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := session.ParseToken(token, w.secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	room := w.lobby.GetRoom(claims["roomId"].(string))
	room.InitClient(
		claims["clientId"].(string),
		claims["username"].(string),
		conn,
	)
}

func (w *Worker) HandleCreate(rw http.ResponseWriter, r *http.Request) {
	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		message := "unable to parse create request"
		http.Error(rw, message, http.StatusBadRequest)
		return
	}

	w.lobby.CreateRoom(request.RoomId)
	rw.WriteHeader(http.StatusCreated)
}

func (w *Worker) HandleStatus(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
