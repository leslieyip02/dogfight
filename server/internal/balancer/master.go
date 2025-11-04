package balancer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"server/internal/id"
	"server/internal/session"
	"server/pb"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/protobuf/proto"
)

const (
	HTTP_TIMEOUT   = 5 * time.Second
	PROBE_INTERVAL = 60 * time.Second
)

func NewRegisterRequest(host string) *pb.RegisterRequest {
	return &pb.RegisterRequest{
		Host: host,
	}
}

type Master struct {
	host string
	port string

	client http.Client
	secret []byte

	roomCapacity        int // max number of clients that can be assigned
	hostOccupancies     map[string]int
	roomOccupancies     map[string]int
	hostToRoomsRegistry map[string][]string // mapping of host to room IDs
	roomToHostRegistry  map[string]string   // mapping of room ID to host

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMaster(
	host string,
	port string,
	secret []byte,
	roomCapacity int,
) *Master {
	client := http.Client{Timeout: HTTP_TIMEOUT}
	ctx, cancel := context.WithCancel(context.Background())

	return &Master{
		host:                host,
		port:                port,
		client:              client,
		secret:              secret,
		roomCapacity:        roomCapacity,
		hostOccupancies:     map[string]int{},
		roomOccupancies:     map[string]int{},
		hostToRoomsRegistry: map[string][]string{},
		roomToHostRegistry:  map[string]string{},
		mu:                  sync.Mutex{},
		ctx:                 ctx,
		cancel:              cancel,
	}
}

func (m *Master) Serve() {
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	r := chi.NewRouter()
	r.Use(corsHandler)
	r.Use(middleware.Logger)

	buildDir := http.Dir("../client/dist")
	fs := http.FileServer(buildDir)
	r.Handle("/*", fs)

	r.Post("/api/join", m.HandleJoin)
	r.Put("/internal/register", m.HandleRegister)

	go m.probeWorkers()

	log.Printf("server is running on http://%s%s", m.host, m.port)
	if err := http.ListenAndServe(m.port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (m *Master) HandleRegister(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request pb.RegisterRequest
	err = proto.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request pb.JoinRequest
	err = proto.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	host, roomId, err := m.getHost(request.RoomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("preparing to send client to %s in room %s", host, roomId)

	m.mu.Lock()
	defer m.mu.Unlock()

	// These occupancies will desync when players leave the room,
	// but will be synced again through periodic status probes
	m.hostOccupancies[host]++
	m.roomOccupancies[roomId]++

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

	body, err := proto.Marshal(&pb.JoinResponse{
		ClientId: clientId,
		Host:     host,
		Token:    token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *Master) getHost(roomId *string) (string, string, error) {
	if roomId != nil {
		return m.lookupHost(*roomId)
	}
	return m.assignHost()

}

func (m *Master) lookupHost(roomId string) (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	host, found := m.roomToHostRegistry[roomId]
	if !found {
		return "", "", fmt.Errorf("room %s not found", roomId)
	}
	if m.roomOccupancies[roomId] >= m.roomCapacity {
		return "", "", fmt.Errorf("room %s is full", roomId)
	}
	return host, roomId, nil
}

func (m *Master) assignHost() (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Look for any available rooms
	for roomId, host := range m.roomToHostRegistry {
		if m.roomOccupancies[roomId] < int(m.roomCapacity) {
			return host, roomId, nil
		}
	}

	// Create a new room if there are no available rooms
	roomId, err := id.NewShortId()
	if err != nil {
		return "", "", err
	}

	host, err := m.chooseHost()
	if err != nil {
		return "", "", err
	}

	err = m.createRoom(host, roomId)
	if err != nil {
		return "", "", err
	}

	m.roomToHostRegistry[roomId] = host
	m.hostToRoomsRegistry[host] = append(m.hostToRoomsRegistry[host], roomId)
	return host, roomId, nil
}

func (m *Master) createRoom(host string, roomId string) error {
	body, err := proto.Marshal(&pb.CreateRequest{
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

func (m *Master) chooseHost() (string, error) {
	// Least connection
	var chosen *string = nil
	least := math.MaxInt
	for host, occupancy := range m.hostOccupancies {
		if occupancy < least {
			chosen = &host
			least = occupancy
		}
	}

	if chosen == nil {
		return "", fmt.Errorf("no hosts available")
	}
	return *chosen, nil
}

func (m *Master) probeWorkers() {
	ticker := time.NewTicker(PROBE_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return

		case <-ticker.C:
			// Probe workers asynchronously
			for host := range m.hostOccupancies {
				go m.probe(host)
			}
		}
	}
}

func (m *Master) probe(host string) {
	url := fmt.Sprintf("http://%s/internal/status", host)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("failed to get status from %s", host)
		return
	}

	response, err := m.client.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("failed to get status from %s", host)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read status from %s", host)
		return
	}

	body := pb.StatusResponse{}
	err = proto.Unmarshal(data, &body)
	if err != nil {
		log.Printf("failed to parse status from %s", host)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Overwrite occupancies with the most recent status
	hostOccupancy := 0
	for _, roomStatus := range body.RoomStatuses {
		m.roomOccupancies[roomStatus.RoomId] = int(roomStatus.Occupancy)
		hostOccupancy += int(roomStatus.Occupancy)
	}
	m.hostOccupancies[host] = hostOccupancy
}
