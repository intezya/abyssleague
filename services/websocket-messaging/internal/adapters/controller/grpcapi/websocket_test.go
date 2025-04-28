package grpcapi

import (
	"context"
	"errors"
	"testing"

	websocketpb "github.com/intezya/abyssleague/proto/websocket"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	errServiceError     = errors.New("service error")
	errUserNotConnected = errors.New("user not connected")
	errBroadcastFailed  = errors.New("broadcast failed")
)

// MockWebsocketService is a mock implementation of the WebsocketService interface.
type MockWebsocketService struct {
	GetOnlineFunc      func(ctx context.Context) (int, error)
	GetOnlineUsersFunc func(ctx context.Context) ([]*service.OnlineUser, error)
	SendToUserFunc     func(ctx context.Context, userID int, jsonPayload []byte) error
	BroadcastFunc      func(ctx context.Context, jsonPayload []byte) error
}

func (m *MockWebsocketService) GetOnline(ctx context.Context) (int, error) {
	return m.GetOnlineFunc(ctx)
}

func (m *MockWebsocketService) GetOnlineUsers(ctx context.Context) ([]*service.OnlineUser, error) {
	return m.GetOnlineUsersFunc(ctx)
}

func (m *MockWebsocketService) SendToUser(
	ctx context.Context,
	userID int,
	jsonPayload []byte,
) error {
	return m.SendToUserFunc(ctx, userID, jsonPayload)
}

func (m *MockWebsocketService) Broadcast(ctx context.Context, jsonPayload []byte) error {
	return m.BroadcastFunc(ctx, jsonPayload)
}

func TestGetOnline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		mockGetOnline  func(ctx context.Context) (int, error)
		expectedResult *websocketpb.GetOnlineResponse
		expectedError  error
	}{
		{
			name: "Success",
			mockGetOnline: func(ctx context.Context) (int, error) {
				return 42, nil
			},
			expectedResult: &websocketpb.GetOnlineResponse{
				Online: 42,
			},
			expectedError: nil,
		},
		{
			name: "Service Error",
			mockGetOnline: func(ctx context.Context) (int, error) {
				return 0, errServiceError
			},
			expectedResult: nil,
			expectedError:  InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				mockService := &MockWebsocketService{
					GetOnlineFunc: tt.mockGetOnline,
				}
				handler := NewWebsocketHandler(mockService)

				result, err := handler.GetOnline(t.Context(), &emptypb.Empty{})

				if tt.expectedError != nil {
					if err == nil {
						t.Errorf("Expected error %v, got nil", tt.expectedError)

						return
					}

					statusErr, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %v", err)

						return
					}

					expectedStatusErr, _ := status.FromError(tt.expectedError)
					if statusErr.Code() != expectedStatusErr.Code() ||
						statusErr.Message() != expectedStatusErr.Message() {
						t.Errorf("Expected error %v, got %v", tt.expectedError, err)
					}

					return
				}

				if err != nil {
					t.Errorf("Unexpected error: %v", err)

					return
				}

				if result.GetOnline() != tt.expectedResult.GetOnline() {
					t.Errorf(
						"Expected online count %d, got %d",
						tt.expectedResult.GetOnline(),
						result.GetOnline(),
					)
				}
			},
		)
	}
}

func TestGetOnlineUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		mockGetOnlineUsers func(ctx context.Context) ([]*service.OnlineUser, error)
		expectedCount      int
		expectedError      error
	}{
		{
			name: "Success",
			mockGetOnlineUsers: func(ctx context.Context) ([]*service.OnlineUser, error) {
				return []*service.OnlineUser{
					{Id: 1, Username: "user1", HardwareID: "hw1"},
					{Id: 2, Username: "user2", HardwareID: "hw2"},
				}, nil
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Empty Result",
			mockGetOnlineUsers: func(ctx context.Context) ([]*service.OnlineUser, error) {
				return []*service.OnlineUser{}, nil
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "Service Error",
			mockGetOnlineUsers: func(ctx context.Context) ([]*service.OnlineUser, error) {
				return nil, errServiceError
			},
			expectedCount: 0,
			expectedError: InternalError,
		},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				testGetOnlineUsersCase(t, tt)
			},
		)
	}
}

func testGetOnlineUsersCase(
	t *testing.T, tt struct {
		name               string
		mockGetOnlineUsers func(ctx context.Context) ([]*service.OnlineUser, error)
		expectedCount      int
		expectedError      error
	},
) {
	t.Helper()

	mockService := &MockWebsocketService{
		GetOnlineUsersFunc: tt.mockGetOnlineUsers,
	}
	handler := NewWebsocketHandler(mockService)

	result, err := handler.GetOnlineUsers(t.Context(), &emptypb.Empty{})

	// Handle expected error case
	if tt.expectedError != nil {
		verifyError(t, err, tt.expectedError)

		return
	}

	// Handle unexpected error case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)

		return
	}

	// Verify result count
	if len(result.GetUsers()) != tt.expectedCount {
		t.Errorf("Expected %d users, got %d", tt.expectedCount, len(result.GetUsers()))

		return
	}

	// Verify first user data if we have results
	if tt.expectedCount > 0 {
		user := result.GetUsers()[0]
		if user.GetId() != 1 || user.GetUsername() != "user1" || user.GetHardwareID() != "hw1" {
			t.Errorf("User data mismatch: %v", user)
		}
	}
}

func verifyError(t *testing.T, err, expectedError error) {
	t.Helper()

	// Check that we got an error
	if err == nil {
		t.Errorf("Expected error %v, got nil", expectedError)

		return
	}

	// Check that it's a gRPC status error
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Errorf("Expected gRPC status error, got %v", err)

		return
	}

	// Check that it's the expected error
	expectedStatusErr, _ := status.FromError(expectedError)
	if statusErr.Code() != expectedStatusErr.Code() ||
		statusErr.Message() != expectedStatusErr.Message() {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		request        *websocketpb.SendMessageRequest
		mockSendToUser func(ctx context.Context, userID int, jsonPayload []byte) error
		expectedError  error
		expectedCode   codes.Code
	}{
		{
			name: "Success",
			request: &websocketpb.SendMessageRequest{
				UserId:      123,
				JsonPayload: []byte(`{"message":"test"}`),
			},
			mockSendToUser: func(ctx context.Context, userID int, jsonPayload []byte) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "Missing UserId",
			request: &websocketpb.SendMessageRequest{
				UserId:      0,
				JsonPayload: []byte(`{"message":"test"}`),
			},
			mockSendToUser: func(ctx context.Context, userID int, jsonPayload []byte) error {
				return nil
			},
			expectedError: status.Errorf(codes.InvalidArgument, "UserId is required"),
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "Missing JsonPayload",
			request: &websocketpb.SendMessageRequest{
				UserId:      123,
				JsonPayload: nil,
			},
			mockSendToUser: func(ctx context.Context, userID int, jsonPayload []byte) error {
				return nil
			},
			expectedError: status.Errorf(codes.InvalidArgument, "JsonPayload is required"),
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "User Not Connected",
			request: &websocketpb.SendMessageRequest{
				UserId:      123,
				JsonPayload: []byte(`{"message":"test"}`),
			},
			mockSendToUser: func(ctx context.Context, userID int, jsonPayload []byte) error {
				return errUserNotConnected
			},
			expectedError: status.Errorf(codes.NotFound, "%s", errUserNotConnected.Error()),
			expectedCode:  codes.NotFound,
		},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				testSendMessageCase(t, tt)
			},
		)
	}
}

func testSendMessageCase(
	t *testing.T, tt struct {
		name           string
		request        *websocketpb.SendMessageRequest
		mockSendToUser func(ctx context.Context, userID int, jsonPayload []byte) error
		expectedError  error
		expectedCode   codes.Code
	},
) {
	t.Helper()

	mockService := &MockWebsocketService{
		SendToUserFunc: tt.mockSendToUser,
	}
	handler := NewWebsocketHandler(mockService)

	_, err := handler.SendMessage(t.Context(), tt.request)

	if tt.expectedError != nil {
		if err == nil {
			t.Errorf("Expected error %v, got nil", tt.expectedError)

			return
		}

		statusErr, ok := status.FromError(err)
		if !ok {
			t.Errorf("Expected gRPC status error, got %v", err)

			return
		}

		if statusErr.Code() != tt.expectedCode {
			t.Errorf("Expected error code %v, got %v", tt.expectedCode, statusErr.Code())
		}

		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBroadcast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		request       *websocketpb.BroadcastRequest
		mockBroadcast func(ctx context.Context, jsonPayload []byte) error
		expectedError error
		expectedCode  codes.Code
	}{
		{
			name: "Success",
			request: &websocketpb.BroadcastRequest{
				JsonPayload: []byte(`{"message":"broadcast"}`),
			},
			mockBroadcast: func(ctx context.Context, jsonPayload []byte) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "Missing JsonPayload",
			request: &websocketpb.BroadcastRequest{
				JsonPayload: nil,
			},
			mockBroadcast: func(ctx context.Context, jsonPayload []byte) error {
				return nil
			},
			expectedError: status.Errorf(codes.InvalidArgument, "JsonPayload is required"),
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "Service Error",
			request: &websocketpb.BroadcastRequest{
				JsonPayload: []byte(`{"message":"broadcast"}`),
			},
			mockBroadcast: func(ctx context.Context, jsonPayload []byte) error {
				return errBroadcastFailed
			},
			expectedError: InternalError,
			expectedCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				mockService := &MockWebsocketService{
					BroadcastFunc: tt.mockBroadcast,
				}
				handler := NewWebsocketHandler(mockService)

				_, err := handler.Broadcast(t.Context(), tt.request)

				if tt.expectedError != nil {
					if err == nil {
						t.Errorf("Expected error %v, got nil", tt.expectedError)

						return
					}

					statusErr, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %v", err)

						return
					}

					if statusErr.Code() != tt.expectedCode {
						t.Errorf("Expected error code %v, got %v", tt.expectedCode, statusErr.Code())
					}

					return
				}

				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			},
		)
	}
}
