package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
	"github.com/intezya/pkglib/logger"
)

func TestNewMiddleware(t *testing.T) {
	t.Parallel()

	// Create a real JWT helper with test configuration
	jwtConfig := auth.NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	jwtHelper := auth.NewJWTHelper(jwtConfig)

	// Create the middleware
	middleware := NewMiddleware(jwtHelper)

	// Verify that the middleware is created correctly
	if middleware == nil {
		t.Fatal("Expected non-nil middleware")
	}

	// Since we've verified middleware is not nil, we can safely check jwtService
	if middleware.jwtService == nil {
		t.Error("Expected non-nil jwtService")
	}
}

func TestJwtAuth(t *testing.T) {
	t.Parallel()

	// Initialize the logger
	_, err := logger.New(logger.WithDebug(true))
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create a real JWT helper with test configuration
	jwtConfig := auth.NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	jwtHelper := auth.NewJWTHelper(jwtConfig)

	// Create a valid token
	validTokenData := &auth.TokenData{
		ID:       123,
		Username: "testuser",
		Hwid:     "testhwid",
	}
	validToken := jwtHelper.TokenGenerator(validTokenData)

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
		expectedAuthData   *entity.AuthenticationData
	}{
		{
			name:               "Valid token",
			authHeader:         "Bearer " + validToken,
			expectedStatusCode: http.StatusOK,
			expectedAuthData:   entity.NewAuthenticationData(123, "testuser", "testhwid"),
		},
		{
			name:               "Missing authorization header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedAuthData:   nil,
		},
		{
			name:               "Invalid token",
			authHeader:         "Bearer invalid-token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedAuthData:   nil,
		},
		{
			name:               "Token with different prefix",
			authHeader:         "Token " + validToken,
			expectedStatusCode: http.StatusOK,
			expectedAuthData:   entity.NewAuthenticationData(123, "testuser", "testhwid"),
		},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create the middleware
			middleware := NewMiddleware(jwtHelper)

			// Create a test request with context
			ctx := t.Context()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add the authorization header if provided
			if tt.authHeader != "" {
				req.Header.Add("Authorization", tt.authHeader)
			}

			// Create a test response recorder
			recorder := httptest.NewRecorder()

			// Call the middleware
			authData := middleware.JwtAuth(recorder, req)

			// Check the status code
			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, recorder.Code)
			}

			// Check the authentication data
			verifyAuthData(t, authData, tt.expectedAuthData)
		})
	}
}

// Helper function to verify authentication data.
func verifyAuthData(t *testing.T, actual, expected *entity.AuthenticationData) {
	t.Helper()

	if expected == nil {
		if actual != nil {
			t.Errorf("Expected nil auth data, got %v", actual)
		}

		return
	}

	if actual == nil {
		t.Error("Expected non-nil auth data, got nil")

		return
	}

	if actual.ID() != expected.ID() {
		t.Errorf("Expected auth data ID %d, got %d", expected.ID(), actual.ID())
	}

	if actual.Username() != expected.Username() {
		t.Errorf("Expected auth data Username %s, got %s", expected.Username(), actual.Username())
	}

	if actual.HardwareID() != expected.HardwareID() {
		t.Errorf(
			"Expected auth data HardwareID %s, got %s",
			expected.HardwareID(),
			actual.HardwareID(),
		)
	}
}

func TestSoftExtractTokenFromHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		tokenString    string
		prefixes       []string
		expectedResult string
	}{
		{
			name:           "Bearer prefix",
			tokenString:    "Bearer token123",
			prefixes:       []string{"Bearer "},
			expectedResult: "token123",
		},
		{
			name:           "Token prefix",
			tokenString:    "Token token123",
			prefixes:       []string{"Bearer ", "Token "},
			expectedResult: "token123",
		},
		{
			name:           "No matching prefix",
			tokenString:    "token123",
			prefixes:       []string{"Bearer ", "Token "},
			expectedResult: "token123",
		},
		{
			name:           "Empty token",
			tokenString:    "",
			prefixes:       []string{"Bearer "},
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := softExtractTokenFromHeader(tt.tokenString, tt.prefixes...)
			if result != tt.expectedResult {
				t.Errorf("Expected result '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}
