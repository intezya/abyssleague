package auth

import (
	"testing"
	"time"
)

func TestJWTConfiguration(t *testing.T) {
	// Test creating a new JWT configuration
	secretKey := "test-secret"
	issuer := "test-issuer"
	expirationTime := 24 * time.Hour

	config := NewJWTConfiguration(secretKey, issuer, expirationTime)

	if config == nil {
		t.Fatal("Expected non-nil JWTConfiguration")
	}
}

func TestNewJWTHelper(t *testing.T) {
	// Test creating a new JWT helper
	config := NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	helper := NewJWTHelper(config)

	if helper == nil {
		t.Fatal("Expected non-nil JWTHelper")
	}

	if helper.JWTConfiguration != config {
		t.Errorf("Expected JWTConfiguration to be %v, got %v", config, helper.JWTConfiguration)
	}
}

func TestTokenGeneratorAndValidation(t *testing.T) {
	// Create a JWT helper with a short expiration time for testing
	config := NewJWTConfiguration("test-secret", "test-issuer", 2*time.Second)
	helper := NewJWTHelper(config)

	// Create test token data
	tokenData := &TokenData{
		ID:       123,
		Username: "testuser",
		Hwid:     "testhwid",
	}

	// Generate a token
	token := helper.TokenGenerator(tokenData)

	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	// Validate the token
	validatedData, err := helper.ValidateToken(token)
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}

	// Check that the validated data matches the original data
	if validatedData.ID != tokenData.ID {
		t.Errorf("Expected ID %d, got %d", tokenData.ID, validatedData.ID)
	}
	if validatedData.Username != tokenData.Username {
		t.Errorf("Expected Username %s, got %s", tokenData.Username, validatedData.Username)
	}
	if validatedData.Hwid != tokenData.Hwid {
		t.Errorf("Expected Hwid %s, got %s", tokenData.Hwid, validatedData.Hwid)
	}

	// Test token expiration
	t.Run("Token expiration", func(t *testing.T) {
		// Wait for the token to expire
		time.Sleep(3 * time.Second)

		// Try to validate the expired token
		_, err := helper.ValidateToken(token)
		if err == nil {
			t.Error("Expected error for expired token, got nil")
		}
	})
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
	// Create a JWT helper
	config := NewJWTConfiguration("test-secret", "test-issuer", 24*time.Hour)
	helper := NewJWTHelper(config)

	// Test cases for invalid tokens
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid format",
			token: "not-a-jwt-token",
		},
		{
			name:  "Wrong signature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := helper.ValidateToken(tc.token)
			if err == nil {
				t.Errorf("Expected error for invalid token, got nil")
			}
		})
	}
}

func TestValidateTokenWithWrongIssuer(t *testing.T) {
	// Create a JWT helper with a specific issuer
	config1 := NewJWTConfiguration("test-secret", "issuer1", 24*time.Hour)
	helper1 := NewJWTHelper(config1)

	// Create another JWT helper with a different issuer
	config2 := NewJWTConfiguration("test-secret", "issuer2", 24*time.Hour)
	helper2 := NewJWTHelper(config2)

	// Create test token data
	tokenData := &TokenData{
		ID:       123,
		Username: "testuser",
		Hwid:     "testhwid",
	}

	// Generate a token with the first helper
	token := helper1.TokenGenerator(tokenData)

	// Try to validate the token with the second helper (different issuer)
	_, err := helper2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for token with wrong issuer, got nil")
	}
}
