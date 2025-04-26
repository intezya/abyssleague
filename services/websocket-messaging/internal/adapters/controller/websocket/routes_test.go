package websocket

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/routes"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/hub"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
)

func TestSetupRoute(t *testing.T) {
	t.Parallel()
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
	ctx := t.Context()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		server.URL+routes.WebsocketPathPrefix+"/test-hub",
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// The request should fail with 400 Bad Request (missing upgrade header)
	// or 401 Unauthorized (missing auth token)
	// Either way, it shouldn't be 404 Not Found
	if resp.StatusCode == http.StatusNotFound {
		t.Errorf("Route not found, expected status code other than 404, got %d", resp.StatusCode)
	}
}
