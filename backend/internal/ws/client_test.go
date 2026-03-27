package ws

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// Helpers

type testEnv struct {
	hub    *Hub
	server *httptest.Server
	conn   *websocket.Conn
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	hub := NewHub(newTestConfig())
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(hub.HandleWebSocket))
	t.Cleanup(server.Close)

	wsURL := "ws" + server.URL[4:] + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	time.Sleep(50 * time.Millisecond)
	return &testEnv{hub: hub, server: server, conn: conn}
}

func (e *testEnv) getClient(t *testing.T) *Client {
	t.Helper()
	for c := range e.hub.clients {
		return c
	}
	t.Fatal("no client registered")
	return nil
}

// Tests

func TestClient_CloseSend_OnlyOnce(t *testing.T) {
	client := &Client{send: make(chan []byte, 10)}
	assert.NotPanics(t, func() {
		client.closeSend()
		client.closeSend()
		client.closeSend()
	})
}

func TestClient_CloseSend_Concurrent(t *testing.T) {
	client := &Client{send: make(chan []byte, 10)}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); client.closeSend() }()
	}
	assert.NotPanics(t, wg.Wait)
}

func TestClient_WritePump_SendsMessages(t *testing.T) {
	env := setupTestEnv(t)
	client := env.getClient(t)

	testMsg := []byte(`{"test": "message"}`)
	client.send <- testMsg

	env.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, message, err := env.conn.ReadMessage()
	require.NoError(t, err)
	assert.Equal(t, testMsg, message)
}

func TestClient_ReadPump_UnregistersOnClose(t *testing.T) {
	env := setupTestEnv(t)
	assert.Len(t, env.hub.clients, 1)

	env.conn.Close()
	time.Sleep(100 * time.Millisecond)

	assert.Len(t, env.hub.clients, 0)
}

func TestClient_ReadPump_HandlesReadError(t *testing.T) {
	env := setupTestEnv(t)
	assert.Len(t, env.hub.clients, 1)

	env.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	env.conn.Close()
	time.Sleep(100 * time.Millisecond)

	assert.Len(t, env.hub.clients, 0)
}

func TestClient_WritePump_ClosesOnSendClose(t *testing.T) {
	env := setupTestEnv(t)
	client := env.getClient(t)

	client.closeSend()
	time.Sleep(50 * time.Millisecond)

	env.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err := env.conn.ReadMessage()
	assert.Error(t, err)
}