package balancer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"server/internal/id"
	"server/internal/session"
	"server/pb"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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

	client          http.Client
	hostOccupancies map[string]uint32
	roomOccupancies map[string]uint32
	roomRegistry    map[string]string // mapping of room ID to host
	secret          []byte

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMaster(host string, port string) *Master {
	godotenv.Load()
	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		log.Fatalf("JWT_SECRET must be set")
	}

	client := http.Client{Timeout: HTTP_TIMEOUT}
	ctx, cancel := context.WithCancel(context.Background())

	return &Master{
		host:            host,
		port:            port,
		client:          client,
		hostOccupancies: map[string]uint32{},
		roomOccupancies: map[string]uint32{},
		roomRegistry:    map[string]string{},
		mu:              sync.Mutex{},
		secret:          []byte(secret),
		ctx:             ctx,
		cancel:          cancel,
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

func (m *Master) getVacantHost() (string, error) {
	// TODO: this is just dummy code that returns the first host
	for host := range m.hostOccupancies {
		return host, nil
	}
	return "", fmt.Errorf("no available hosts")
}

func (m *Master) probeWorkers() {
	ticker := time.NewTicker(PROBE_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return

		case <-ticker.C:
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

	var hostOccupancy uint32 = 0
	for _, roomStatus := range body.RoomStatuses {
		m.roomOccupancies[roomStatus.RoomId] = roomStatus.Occupancy
		hostOccupancy += roomStatus.Occupancy
		log.Printf("room %s has %d players", roomStatus.RoomId, roomStatus.Occupancy)
	}
	m.hostOccupancies[host] = hostOccupancy
}
