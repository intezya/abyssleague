package grpcapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalError = status.Errorf(codes.Internal, "An unexpected error occured")
)
