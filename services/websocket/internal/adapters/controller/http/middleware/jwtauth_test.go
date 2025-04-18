package middleware

import (
	"github.com/intezya/pkglib/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"websocket/internal/domain/entity"
	"websocket/internal/pkg/auth"
)

func TestNewMiddleware(t *testing.T) {
	// Create a real JWT helper with test configuration
	jwtConfig := auth.NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	jwtHelper := auth.NewJWTHelper(jwtConfig)

	// Create the middleware
	middleware := NewMiddleware(jwtHelper)

	// Verify that the middleware is not nil
	if middleware == nil {
		t.Error("Expected non-nil middleware")
	}

	// Verify that the middleware has the correct JWT service
	if middleware.jwtService == nil {
		t.Error("Expected non-nil jwtService")
	}
}

func TestJwtAuth(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			// Create the middleware
			middleware := NewMiddleware(jwtHelper)

			// Create a test request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add the authorization header if provided
			if tt.authHeader != "" {
				req.Header.Add("Authorization", tt.authHeader)
			}

			// Create a test response recorder
			rr := httptest.NewRecorder()

			// Call the middleware
			authData := middleware.JwtAuth(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check the authentication data
			if tt.expectedAuthData == nil {
				if authData != nil {
					t.Errorf("Expected nil auth data, got %v", authData)
				}
			} else {
				if authData == nil {
					t.Error("Expected non-nil auth data, got nil")
				} else {
					if authData.ID() != tt.expectedAuthData.ID() {
						t.Errorf("Expected auth data ID %d, got %d", tt.expectedAuthData.ID(), authData.ID())
					}
					if authData.Username() != tt.expectedAuthData.Username() {
						t.Errorf("Expected auth data Username %s, got %s", tt.expectedAuthData.Username(), authData.Username())
					}
					if authData.HardwareID() != tt.expectedAuthData.HardwareID() {
						t.Errorf("Expected auth data HardwareID %s, got %s", tt.expectedAuthData.HardwareID(), authData.HardwareID())
					}
				}
			}
		})
	}
}

func TestSoftExtractTokenFromHeader(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			result := softExtractTokenFromHeader(tt.tokenString, tt.prefixes...)
			if result != tt.expectedResult {
				t.Errorf("Expected result '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}
