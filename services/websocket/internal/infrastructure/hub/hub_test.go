package hub

import (
	"context"
	"github.com/intezya/abyssleague/services/websocket/internal/domain/entity"
	logger2 "github.com/intezya/pkglib/logger"
	"testing"
	"time"
)

// MockClient is a simplified client for testing
type MockClient struct {
	ID          int
	Username    string
	HardwareID  string
	SendChannel chan []byte
}

func (m *MockClient) GetAuthentication() *entity.AuthenticationData {
	return entity.NewAuthenticationData(m.ID, m.Username, m.HardwareID)
}

func TestNewHub(t *testing.T) {
	// Test creating a new hub
	hubName := "test-hub"
	h := NewHub(hubName)

	if h == nil {
		t.Fatal("Expected non-nil Hub")
	}

	if h.name != hubName {
		t.Errorf("Expected hub name %s, got %s", hubName, h.name)
	}

	if h.clients == nil {
		t.Error("Expected non-nil clients map")
	}

	if h.clientsByID == nil {
		t.Error("Expected non-nil clientsByID map")
	}

	if h.register == nil {
		t.Error("Expected non-nil register channel")
	}

	if h.unregister == nil {
		t.Error("Expected non-nil unregister channel")
	}

	if h.broadcast == nil {
		t.Error("Expected non-nil broadcast channel")
	}

	if h.done == nil {
		t.Error("Expected non-nil done channel")
	}
}

func TestGetName(t *testing.T) {
	// Test getting the hub name
	hubName := "test-hub"
	h := NewHub(hubName)

	if h.GetName() != hubName {
		t.Errorf("Expected hub name %s, got %s", hubName, h.GetName())
	}
}

func TestRegisterAndGetClients(t *testing.T) {
	// Create a hub
	_, _ = logger2.New()
	h := NewHub("test-hub")

	// Start the hub
	go h.Run()
	defer h.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{
		Hub:            h,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	h.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Get the clients
	clients := h.GetClients(context.Background())

	// Verify that the client was registered
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(clients))
	}

	if clients[0].ID() != authData.ID() {
		t.Errorf("Expected client ID %d, got %d", authData.ID(), clients[0].ID())
	}

	if clients[0].Username() != authData.Username() {
		t.Errorf("Expected client username %s, got %s", authData.Username(), clients[0].Username())
	}

	if clients[0].HardwareID() != authData.HardwareID() {
		t.Errorf("Expected client hardware ID %s, got %s", authData.HardwareID(), clients[0].HardwareID())
	}
}

func TestUnregisterClient(t *testing.T) {
	// Create a hub
	h := NewHub("test-hub")

	// Start the hub
	go h.Run()
	defer h.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{
		Hub:            h,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	h.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was registered
	clients := h.GetClients(context.Background())
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(clients))
	}

	// Unregister the client
	h.UnregisterClient(client)

	// Give the hub time to process the unregistration
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was unregistered
	clients = h.GetClients(context.Background())
	if len(clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(clients))
	}
}

func TestSendToUser(t *testing.T) {
	// Create a hub
	h := NewHub("test-hub")

	// Start the hub
	go h.Run()
	defer h.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{
		Hub:            h,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	h.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Send a message to the user
	message := []byte("test message")
	success := h.SendToUser(context.Background(), authData.ID(), message)

	// Verify that the message was sent successfully
	if !success {
		t.Error("Expected SendToUser to return true")
	}

	// Verify that the message was received by the client
	select {
	case receivedMessage := <-client.Send:
		if string(receivedMessage) != string(message) {
			t.Errorf("Expected message %s, got %s", string(message), string(receivedMessage))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for message")
	}

	// Test sending to a non-existent user
	success = h.SendToUser(context.Background(), 999, message)
	if success {
		t.Error("Expected SendToUser to return false for non-existent user")
	}
}

func TestBroadcast(t *testing.T) {
	// Create a hub
	h := NewHub("test-hub")

	// Start the hub
	go h.Run()
	defer h.Stop()

	// Create test clients
	client1 := &Client{
		Hub:            h,
		authentication: entity.NewAuthenticationData(1, "user1", "hw1"),
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	client2 := &Client{
		Hub:            h,
		authentication: entity.NewAuthenticationData(2, "user2", "hw2"),
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the clients
	h.RegisterClient(client1)
	h.RegisterClient(client2)

	// Give the hub time to process the registrations
	time.Sleep(100 * time.Millisecond)

	// Broadcast a message
	message := []byte("broadcast message")
	h.Broadcast(context.Background(), message)

	// Give the hub time to process the broadcast
	time.Sleep(100 * time.Millisecond)

	// Verify that both clients received the message
	select {
	case receivedMessage := <-client1.Send:
		if string(receivedMessage) != string(message) {
			t.Errorf("Client 1: Expected message %s, got %s", string(message), string(receivedMessage))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 1: Timed out waiting for message")
	}

	select {
	case receivedMessage := <-client2.Send:
		if string(receivedMessage) != string(message) {
			t.Errorf("Client 2: Expected message %s, got %s", string(message), string(receivedMessage))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 2: Timed out waiting for message")
	}
}
