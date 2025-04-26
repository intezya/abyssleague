package hub

import (
	"testing"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	logger2 "github.com/intezya/pkglib/logger"
)

// MockClient is a simplified client for testing.
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
	t.Parallel()
	// Test creating a new hub
	hubName := "test-hub"
	hub := NewHub(hubName)

	if hub == nil {
		t.Fatal("Expected non-nil Hub")
	}

	if hub.name != hubName {
		t.Errorf("Expected hub name %s, got %s", hubName, hub.name)
	}

	if hub.clients == nil {
		t.Error("Expected non-nil clients map")
	}

	if hub.clientsByID == nil {
		t.Error("Expected non-nil clientsByID map")
	}

	if hub.register == nil {
		t.Error("Expected non-nil register channel")
	}

	if hub.unregister == nil {
		t.Error("Expected non-nil unregister channel")
	}

	if hub.broadcast == nil {
		t.Error("Expected non-nil broadcast channel")
	}

	if hub.done == nil {
		t.Error("Expected non-nil done channel")
	}
}

func TestGetName(t *testing.T) {
	t.Parallel()
	// Test getting the hub name
	hubName := "test-hub"
	hub := NewHub(hubName)

	if hub.GetName() != hubName {
		t.Errorf("Expected hub name %s, got %s", hubName, hub.GetName())
	}
}

func TestRegisterAndGetClients(t *testing.T) {
	t.Parallel()
	// Create a hub
	_, _ = logger2.New()
	hub := NewHub("test-hub")

	// Start the hub
	go hub.Run()
	defer hub.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{ //nolint:exhaustruct
		Hub:            hub,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	hub.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Get the clients
	clients := hub.GetClients(t.Context())

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
		t.Errorf(
			"Expected client hardware ID %s, got %s",
			authData.HardwareID(),
			clients[0].HardwareID(),
		)
	}
}

func TestUnregisterClient(t *testing.T) {
	t.Parallel()
	// Create a hub
	hub := NewHub("test-hub")

	// Start the hub
	go hub.Run()
	defer hub.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{ //nolint:exhaustruct
		Hub:            hub,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	hub.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was registered
	clients := hub.GetClients(t.Context())
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(clients))
	}

	// Unregister the client
	hub.UnregisterClient(client)

	// Give the hub time to process the unregistration
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was unregistered
	clients = hub.GetClients(t.Context())
	if len(clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(clients))
	}
}

func TestSendToUser(t *testing.T) {
	t.Parallel()
	// Create a hub
	hub := NewHub("test-hub")

	// Start the hub
	go hub.Run()
	defer hub.Stop()

	// Create a test client
	authData := entity.NewAuthenticationData(1, "testuser", "testhwid")
	client := &Client{ //nolint:exhaustruct
		Hub:            hub,
		authentication: authData,
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the client
	hub.RegisterClient(client)

	// Give the hub time to process the registration
	time.Sleep(100 * time.Millisecond)

	// Send a message to the user
	message := []byte("test message")
	success := hub.SendToUser(t.Context(), authData.ID(), message)

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
	success = hub.SendToUser(t.Context(), 999, message)
	if success {
		t.Error("Expected SendToUser to return false for non-existent user")
	}
}

func TestBroadcast(t *testing.T) {
	t.Parallel()
	// Create a hub
	hub := NewHub("test-hub")

	// Start the hub
	go hub.Run()
	defer hub.Stop()

	// Create test clients
	client1 := &Client{
		Hub:            hub,
		authentication: entity.NewAuthenticationData(1, "user1", "hw1"),
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	client2 := &Client{
		Hub:            hub,
		authentication: entity.NewAuthenticationData(2, "user2", "hw2"),
		Send:           make(chan []byte, 256),
		connectTime:    time.Now(),
	}

	// Register the clients
	hub.RegisterClient(client1)
	hub.RegisterClient(client2)

	// Give the hub time to process the registrations
	time.Sleep(100 * time.Millisecond)

	// Broadcast a message
	message := []byte("broadcast message")
	hub.Broadcast(t.Context(), message)

	// Give the hub time to process the broadcast
	time.Sleep(100 * time.Millisecond)

	// Verify that both clients received the message
	select {
	case receivedMessage := <-client1.Send:
		if string(receivedMessage) != string(message) {
			t.Errorf(
				"Client 1: Expected message %s, got %s",
				string(message),
				string(receivedMessage),
			)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 1: Timed out waiting for message")
	}

	select {
	case receivedMessage := <-client2.Send:
		if string(receivedMessage) != string(message) {
			t.Errorf(
				"Client 2: Expected message %s, got %s",
				string(message),
				string(receivedMessage),
			)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 2: Timed out waiting for message")
	}
}
