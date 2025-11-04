package balancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/internal/id"
	"server/internal/session"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func NewRegisterRequest(host string) *RegisterRequest {
	return &RegisterRequest{
		Host: host,
	}
}

type Master struct {
	url  string
	host string
	port string

	client          http.Client
	hostOccupancies map[string]int
	roomOccupancies map[string]int
	roomRegistry    map[string]string // mapping of room ID to host
	mu              sync.Mutex
	secret          []byte
}

func NewMaster(host string, port string) *Master {
	godotenv.Load()
	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		log.Fatalf("JWT_SECRET must be set")
	}

	return &Master{
		url:             fmt.Sprintf("%s%s", host, port),
		host:            host,
		port:            port,
		client:          http.Client{},
		hostOccupancies: map[string]int{},
		roomOccupancies: map[string]int{},
		roomRegistry:    map[string]string{},
		mu:              sync.Mutex{},
		secret:          []byte(secret),
	}
}

func (m *Master) Serve() {
	r := chi.NewRouter()
	r.Use(corsHandler)
	r.Use(middleware.Logger)

	buildDir := http.Dir("../client/dist")
	fs := http.FileServer(buildDir)
	r.Handle("/*", fs)

	r.Post("/api/join", m.HandleJoin)
	r.Put("/internal/register", m.HandleRegister)

	log.Printf("server is running on %s", m.url)
	if err := http.ListenAndServe(m.port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (m *Master) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		message := "unable to parse register request"
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	host := request.Host + request.Port

	m.mu.Lock()
	defer m.mu.Unlock()

	m.hostOccupancies[host] = 0
	log.Printf("registered %s", host)
	w.WriteHeader(http.StatusCreated)
}

func (m *Master) HandleJoin(w http.ResponseWriter, r *http.Request) {
	var request JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "unable to parse room join request", http.StatusBadRequest)
		return
	}

	host, roomId, err := m.getHost(request.RoomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// TODO: this should be reported by the workers
	m.mu.Lock()
	m.hostOccupancies[host]++
	m.roomOccupancies[roomId]++
	m.mu.Unlock()

	clientId, err := id.NewShortId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := session.CreateToken(clientId, request.Username, roomId, m.secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := JoinResponse{
		ClientId: clientId,
		Host:     host,
		Token:    token,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		message := "unable to write body"
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
}

func (m *Master) getHost(roomId *string) (string, string, error) {
	if roomId == nil {
		return m.assignHost()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	host, found := m.roomRegistry[*roomId]
	if !found {
		return "", "", fmt.Errorf("room %s not found", *roomId)
	}
	return host, *roomId, nil
}

func (m *Master) assignHost() (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("checking roomRegistry: %v", m.roomRegistry)
	for roomId, host := range m.roomRegistry {
		if m.roomOccupancies[roomId] < 32 {
			return host, roomId, nil
		}
	}

	roomId, err := id.NewShortId()
	if err != nil {
		return "", "", err
	}

	host, err := m.getVacantHost()
	if err != nil {
		return "", "", err
	}

	err = m.createRoom(host, roomId)
	if err != nil {
		return "", "", err
	}
	m.roomRegistry[roomId] = host
	return host, roomId, nil
}

func (m *Master) createRoom(host string, roomId string) error {
	body, err := json.Marshal(CreateRequest{
		RoomId: roomId,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/internal/create", host)
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	response, err := m.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %v", response.StatusCode)
	}
	return nil
}

func (m *Master) getVacantHost() (string, error) {
	// TODO: this is just dummy code that returns the first host
	for host := range m.hostOccupancies {
		return host, nil
	}
	return "", fmt.Errorf("no available hosts")
}

// func (m *Master) probeWorkers() {
// 	// TODO: implement
// }

// func (m *Master) getWorkerStatus(url string) error {
// 	// TODO: integrate
// 	request, err := http.NewRequest(
// 		"GET",
// 		fmt.Sprintf("%s/status", url),
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	response, err := m.client.Do(request)
// 	if err != nil {
// 		return err
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		return fmt.Errorf("unexpected status code %v", response.StatusCode)
// 	}

// 	defer response.Body.Close()

// 	body := &StatusResponse{}
// 	err = json.NewDecoder(response.Body).Decode(body)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("%s has %d players", url, body.PlayerCount)

// 	return nil
// }
