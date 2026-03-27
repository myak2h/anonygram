package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"anonygram/internal/config"
	"anonygram/internal/models"
)

// Helpers

func newTestConfig() *config.Config {
	return &config.Config{
		AllowedOrigins:   []string{"*"},
		ClientBufferSize: 256,
		HubBufferSize:    16,
	}
}

func startHub(t *testing.T) *Hub {
	t.Helper()
	hub := NewHub(newTestConfig())
	go hub.Run()
	return hub
}

func newMockClient(hub *Hub, bufSize int) *Client {
	return &Client{hub: hub, send: make(chan []byte, bufSize)}
}

func registerClient(t *testing.T, hub *Hub, bufSize int) *Client {
	t.Helper()
	c := newMockClient(hub, bufSize)
	hub.register <- c
	time.Sleep(10 * time.Millisecond)
	return c
}

// Tests

func TestNewHub(t *testing.T) {
	hub := NewHub(newTestConfig())
	require.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.Equal(t, 256, hub.clientBufferSize)
}

func TestNewUpgrader_CheckOrigin(t *testing.T) {
	tests := []struct {
		origins []string
		origin  string
		allowed bool
	}{
		{[]string{"*"}, "http://any.com", true},
		{[]string{"http://a.com", "http://b.com"}, "http://a.com", true},
		{[]string{"http://a.com"}, "http://evil.com", false},
		{[]string{"http://a.com"}, "", false},
	}
	for _, tt := range tests {
		upgrader := newUpgrader(tt.origins)
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		req.Header.Set("Origin", tt.origin)
		assert.Equal(t, tt.allowed, upgrader.CheckOrigin(req), "origin=%s", tt.origin)
	}
}

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := startHub(t)
	client := registerClient(t, hub, 256)
	assert.Len(t, hub.clients, 1)

	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, hub.clients, 0)
}

func TestHub_Broadcast(t *testing.T) {
	hub := startHub(t)
	assert.True(t, hub.Broadcast(models.Image{ID: "1"}))

	hubFull := NewHub(newTestConfig())
	for i := 0; i < 16; i++ {
		hubFull.Broadcast(models.Image{ID: "fill"})
	}
	assert.False(t, hubFull.Broadcast(models.Image{ID: "overflow"}))
}

func TestHub_BroadcastToSlowClient(t *testing.T) {
	cfg := newTestConfig()
	cfg.ClientBufferSize = 1
	hub := NewHub(cfg)
	go hub.Run()

	slowClient := newMockClient(hub, 1)
	hub.register <- slowClient
	time.Sleep(10 * time.Millisecond)

	hub.Broadcast(models.Image{ID: "msg1"})
	time.Sleep(10 * time.Millisecond)
	hub.Broadcast(models.Image{ID: "msg2"})
	time.Sleep(50 * time.Millisecond)

	assert.Len(t, hub.clients, 0)
}

func TestHub_ConcurrentOperations(t *testing.T) {
	hub := startHub(t)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := newMockClient(hub, 256)
			hub.register <- c
			time.Sleep(5 * time.Millisecond)
			hub.unregister <- c
		}()
	}
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, hub.clients, 0)

	client := registerClient(t, hub, 1000)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); hub.Broadcast(models.Image{ID: "x"}) }()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	received := 0
	for len(client.send) > 0 {
		<-client.send
		received++
	}
	assert.Greater(t, received, 0)
}

func TestHub_HandleWebSocket_Integration(t *testing.T) {
	hub := startHub(t)
	server := httptest.NewServer(http.HandlerFunc(hub.HandleWebSocket))
	t.Cleanup(server.Close)

	wsURL := "ws" + server.URL[4:] + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	time.Sleep(50 * time.Millisecond)
	assert.Len(t, hub.clients, 1)

	hub.Broadcast(models.Image{ID: "ws-test", Title: "WebSocket Test"})

	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, message, err := conn.ReadMessage()
	require.NoError(t, err)

	var img models.Image
	require.NoError(t, json.Unmarshal(message, &img))
	assert.Equal(t, "ws-test", img.ID)
}
