package grpcwrap

var (
	ErrServiceNotAvailable = newError("service not available")
	ErrConnectionFailed    = newError("connection failed")
	ErrRPCTimeout          = newError("rpc request timeout exceeded")
)

type clientError struct {
	message string
}

func (e clientError) Error() string {
	return e.message
}

func newError(message string) error {
	return clientError{message: message}
}
