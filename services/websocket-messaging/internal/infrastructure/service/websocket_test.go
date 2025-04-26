package service

import (
	"context"
	"errors"
	"testing"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
)

var errFailedToSendMessage = errors.New("failed to send message to user")

// MockHub is a mock implementation of the Hub interface for testing.
type MockHub struct {
	GetClientsFunc func(ctx context.Context) []*entity.AuthenticationData
	SendToUserFunc func(ctx context.Context, userId int, jsonPayload []byte) bool
	BroadcastFunc  func(ctx context.Context, jsonPayload []byte)
}

func (m *MockHub) GetClients(ctx context.Context) []*entity.AuthenticationData {
	return m.GetClientsFunc(ctx)
}

func (m *MockHub) SendToUser(ctx context.Context, userId int, jsonPayload []byte) bool {
	return m.SendToUserFunc(ctx, userId, jsonPayload)
}

func (m *MockHub) Broadcast(ctx context.Context, jsonPayload []byte) {
	m.BroadcastFunc(ctx, jsonPayload)
}

// Create a test-specific version of NewWebsocketService that accepts our MockHub.
func newTestWebsocketService(mockHub *MockHub) *WebsocketService {
	return &WebsocketService{hub: mockHub}
}

func TestGetOnline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		mockGetClients func(ctx context.Context) []*entity.AuthenticationData
		expectedCount  int
		expectedError  error
	}{
		{
			name: "Success with clients",
			mockGetClients: func(ctx context.Context) []*entity.AuthenticationData {
				return []*entity.AuthenticationData{
					entity.NewAuthenticationData(1, "user1", "hw1"),
					entity.NewAuthenticationData(2, "user2", "hw2"),
				}
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Success with no clients",
			mockGetClients: func(ctx context.Context) []*entity.AuthenticationData {
				return []*entity.AuthenticationData{}
			},
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHub := &MockHub{ //nolint:exhaustruct
				GetClientsFunc: tt.mockGetClients,
			}
			service := newTestWebsocketService(mockHub)

			count, err := service.GetOnline(t.Context())

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestGetOnlineUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		mockGetClients func(ctx context.Context) []*entity.AuthenticationData
		expectedCount  int
		expectedError  error
	}{
		{
			name: "Success with clients",
			mockGetClients: func(ctx context.Context) []*entity.AuthenticationData {
				return []*entity.AuthenticationData{
					entity.NewAuthenticationData(1, "user1", "hw1"),
					entity.NewAuthenticationData(2, "user2", "hw2"),
				}
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Success with no clients",
			mockGetClients: func(ctx context.Context) []*entity.AuthenticationData {
				return []*entity.AuthenticationData{}
			},
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHub := &MockHub{ //nolint:exhaustruct
				GetClientsFunc: tt.mockGetClients,
			}
			service := newTestWebsocketService(mockHub)

			users, err := service.GetOnlineUsers(t.Context())

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if len(users) != tt.expectedCount {
				t.Errorf("Expected %d users, got %d", tt.expectedCount, len(users))

				return
			}

			if tt.expectedCount > 0 {
				// Verify first user data
				if users[0].Id != 1 || users[0].Username != "user1" ||
					users[0].HardwareID != "hw1" {
					t.Errorf("User data mismatch: %v", users[0])
				}
			}
		})
	}
}

func TestSendToUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		userId         int
		jsonPayload    []byte
		mockSendToUser func(ctx context.Context, userId int, jsonPayload []byte) bool
		expectedError  error
	}{
		{
			name:        "Success",
			userId:      1,
			jsonPayload: []byte(`{"message":"test"}`),
			mockSendToUser: func(ctx context.Context, userId int, jsonPayload []byte) bool {
				return true
			},
			expectedError: nil,
		},
		{
			name:        "User not found",
			userId:      999,
			jsonPayload: []byte(`{"message":"test"}`),
			mockSendToUser: func(ctx context.Context, userId int, jsonPayload []byte) bool {
				return false
			},
			expectedError: errFailedToSendMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHub := &MockHub{ //nolint:exhaustruct
				SendToUserFunc: tt.mockSendToUser,
			}
			service := newTestWebsocketService(mockHub)

			err := service.SendToUser(t.Context(), tt.userId, tt.jsonPayload)

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)

				return
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error message %q, got %q", tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestBroadcast(t *testing.T) {
	t.Parallel()

	broadcastCalled := false
	mockHub := &MockHub{ //nolint:exhaustruct
		BroadcastFunc: func(ctx context.Context, jsonPayload []byte) {
			broadcastCalled = true
		},
	}
	service := newTestWebsocketService(mockHub)

	err := service.Broadcast(t.Context(), []byte(`{"message":"broadcast"}`))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !broadcastCalled {
		t.Errorf("Expected Broadcast to be called")
	}
}
