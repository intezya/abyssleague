package entity

import (
	"reflect"
	"strconv"
	"testing"
	"websocket/internal/pkg/auth"
)

func TestNewAuthenticationData(t *testing.T) {
	// Test creating a new authentication data object
	id := 123
	username := "testuser"
	hwid := "testhwid"

	authData := NewAuthenticationData(id, username, hwid)

	if authData == nil {
		t.Fatal("Expected non-nil AuthenticationData")
	}

	// Test getters
	if authData.ID() != id {
		t.Errorf("Expected ID %d, got %d", id, authData.ID())
	}

	if authData.Username() != username {
		t.Errorf("Expected Username %s, got %s", username, authData.Username())
	}

	if authData.HardwareID() != hwid {
		t.Errorf("Expected HardwareID %s, got %s", hwid, authData.HardwareID())
	}
}

func TestAuthenticationData_Encode(t *testing.T) {
	// Test encoding to a map
	authData := NewAuthenticationData(123, "testuser", "testhwid")

	encoded := authData.Encode()

	expectedMap := map[string]string{
		"id":       "123",
		"username": "testuser",
		"hwid":     "testhwid",
	}

	if !reflect.DeepEqual(encoded, expectedMap) {
		t.Errorf("Expected encoded map %v, got %v", expectedMap, encoded)
	}
}

func TestDecodeToAuthenticationData(t *testing.T) {
	// Test decoding from a TokenData
	tokenData := &auth.TokenData{
		ID:       123,
		Username: "testuser",
		Hwid:     "testhwid",
	}

	authData := DecodeToAuthenticationData(tokenData)

	if authData == nil {
		t.Fatal("Expected non-nil AuthenticationData")
	}

	if authData.ID() != tokenData.ID {
		t.Errorf("Expected ID %d, got %d", tokenData.ID, authData.ID())
	}

	if authData.Username() != tokenData.Username {
		t.Errorf("Expected Username %s, got %s", tokenData.Username, authData.Username())
	}

	if authData.HardwareID() != tokenData.Hwid {
		t.Errorf("Expected HardwareID %s, got %s", tokenData.Hwid, authData.HardwareID())
	}
}

func TestAuthenticationData_Getters(t *testing.T) {
	// Test all getters with different values
	testCases := []struct {
		id       int
		username string
		hwid     string
	}{
		{1, "user1", "hw1"},
		{2, "user2", "hw2"},
		{0, "", ""},
		{-1, "negative", "negative-hw"},
	}

	for _, tc := range testCases {
		t.Run("ID:"+strconv.Itoa(tc.id), func(t *testing.T) {
			authData := NewAuthenticationData(tc.id, tc.username, tc.hwid)

			if authData.ID() != tc.id {
				t.Errorf("Expected ID %d, got %d", tc.id, authData.ID())
			}

			if authData.Username() != tc.username {
				t.Errorf("Expected Username %s, got %s", tc.username, authData.Username())
			}

			if authData.HardwareID() != tc.hwid {
				t.Errorf("Expected HardwareID %s, got %s", tc.hwid, authData.HardwareID())
			}
		})
	}
}
