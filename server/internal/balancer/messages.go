package balancer

type RegisterRequest struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RegisterResponse struct{}

type JoinRequest struct {
	Username string  `json:"username"`
	RoomId   *string `json:"roomId,omitempty"`
}

type JoinResponse struct {
	ClientId string `json:"clientId"`
	Host     string `json:"host"`
	Token    string `json:"token"`
}

type StatusResponse struct {
	PlayerCount int      `json:"playerCount"`
	RoomIds     []string `json:"roomIds"`
}

type CreateRequest struct {
	RoomId string `json:"roomId"`
}
