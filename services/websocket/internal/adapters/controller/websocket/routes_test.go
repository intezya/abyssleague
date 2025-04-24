package websocket

import (
	"github.com/intezya/abyssleague/services/websocket/internal/adapters/controller/http/routes"
	"github.com/intezya/abyssleague/services/websocket/internal/infrastructure/hub"
	"github.com/intezya/abyssleague/services/websocket/internal/pkg/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetupRoute(t *testing.T) {
	// Create a test HTTP mux
	mux := http.NewServeMux()

	// Create a test hub
	testHub := hub.NewHub("test-hub")

	// Create a test JWT helper
	jwtConfig := auth.NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	jwtHelper := auth.NewJWTHelper(jwtConfig)

	// Call the function being tested
	SetupRoute(mux, testHub, "test-hub", jwtHelper)

	// Create a test server using the mux
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test that the route was registered by making a request
	// Note: We don't expect a successful websocket connection here,
	// just verifying that the route handler was registered
	resp, err := http.Get(server.URL + routes.WebsocketPathPrefix + "/test-hub")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// The request should fail with 400 Bad Request (missing upgrade header)
	// or 401 Unauthorized (missing auth token)
	// Either way, it shouldn't be 404 Not Found
	if resp.StatusCode == http.StatusNotFound {
		t.Errorf("Route not found, expected status code other than 404, got %d", resp.StatusCode)
	}
}
