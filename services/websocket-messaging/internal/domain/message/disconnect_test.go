package message

import (
	"encoding/json"
	"testing"
)

func TestDisconnectByOtherClient(t *testing.T) {
	t.Parallel()
	// Parse the DisconnectByOtherClient message
	var msg message

	err := json.Unmarshal(DisconnectByOtherClient, &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal DisconnectByOtherClient: %v", err)
	}

	// Verify the message fields
	if msg.Type != "disconnect" {
		t.Errorf("Expected Type to be 'disconnect', got '%s'", msg.Type)
	}

	if msg.Subtype != "other_client" {
		t.Errorf("Expected Subtype to be 'other_client', got '%s'", msg.Subtype)
	}

	if msg.Message != "You have been disconnected by another connection" {
		t.Errorf(
			"Expected message to be 'You have been disconnected by another connection', got '%s'",
			msg.Message,
		)
	}
}

func TestMessageStruct(t *testing.T) {
	t.Parallel()
	// Create a message struct
	msg := message{
		Type:    "test_type",
		Subtype: "test_subtype",
		Message: "test_message",
	}

	// Marshal the message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Unmarshal the JSON back to a message struct
	var unmarshaledMsg message

	err = json.Unmarshal(data, &unmarshaledMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify the unmarshaled message fields
	if unmarshaledMsg.Type != "test_type" {
		t.Errorf("Expected Type to be 'test_type', got '%s'", unmarshaledMsg.Type)
	}

	if unmarshaledMsg.Subtype != "test_subtype" {
		t.Errorf("Expected Subtype to be 'test_subtype', got '%s'", unmarshaledMsg.Subtype)
	}

	if unmarshaledMsg.Message != "test_message" {
		t.Errorf("Expected message to be 'test_message', got '%s'", unmarshaledMsg.Message)
	}
}
