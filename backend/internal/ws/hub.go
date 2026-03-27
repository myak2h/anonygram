package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"anonygram/internal/config"
	"anonygram/internal/models"
)

type Broadcaster interface {
	Broadcast(img models.Image) bool
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
}

type Hub struct {
	clients          map[*Client]bool
	broadcast        chan models.Image
	register         chan *Client
	unregister       chan *Client
	upgrader         websocket.Upgrader
	clientBufferSize int
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		clients:          make(map[*Client]bool),
		broadcast:        make(chan models.Image, cfg.HubBufferSize),
		register:         make(chan *Client, cfg.HubBufferSize),
		unregister:       make(chan *Client, cfg.HubBufferSize),
		upgrader:         newUpgrader(cfg.AllowedOrigins),
		clientBufferSize: cfg.ClientBufferSize,
	}
}

func newUpgrader(allowedOrigins []string) websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.closeSend()
			}
		case img := <-h.broadcast:
			data, err := json.Marshal(img)
			if err != nil {
				log.Printf("failed to marshal image for websocket broadcast: %v", err)
				continue
			}
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					client.closeSend()
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, h.clientBufferSize)}
	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) Broadcast(img models.Image) bool {
	select {
	case h.broadcast <- img:
		return true
	default:
		return false
	}
}
