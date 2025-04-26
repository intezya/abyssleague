package grpcapi

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInternalError(t *testing.T) {
	t.Parallel()
	// Verify that InternalError is a gRPC status error with the Internal code
	statusErr, ok := status.FromError(InternalError)
	if !ok {
		t.Errorf("InternalError is not a gRPC status error")
	}

	if statusErr.Code() != codes.Internal {
		t.Errorf("Expected error code %v, got %v", codes.Internal, statusErr.Code())
	}

	if statusErr.Message() != "An unexpected error occurred" {
		t.Errorf(
			"Expected error message 'An unexpected error occured', got '%v'",
			statusErr.Message(),
		)
	}
}
