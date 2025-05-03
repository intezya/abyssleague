package grpcwrap

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	connectionStateDisconnected int32 = iota
	connectionStateConnecting
	connectionStateConnected
	connectionStateFailed

	defaultRpcTimeout        = 2 * time.Second
	defaultConnectionTimeout = 5 * time.Second
)

type ClientOption func(client *BaseGRPCClient)

// WithDevMode enables dev mode, that ignores empty client error.
func WithDevMode(devMode bool) ClientOption {
	return func(client *BaseGRPCClient) {
		client.DevMode = devMode
	}
}

// WithDialOptions adds additional options for gRPC connection.
func WithDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(client *BaseGRPCClient) {
		client.dialOptions = append(client.dialOptions, opts...)
	}
}

// WithConnectionTimeout sets connection timeout.
func WithConnectionTimeout(timeout time.Duration) ClientOption {
	return func(client *BaseGRPCClient) {
		client.connectionTimeout = timeout
	}
}

// WithRPCTimeout setup timeout for RPC calls.
func WithRPCTimeout(timeout time.Duration) ClientOption {
	return func(client *BaseGRPCClient) {
		client.rpcTimeout = timeout
	}
}

// ClientCreator - function for gRPC client creation from connection.
type ClientCreator[T any] func(conn *grpc.ClientConn) T

type TypeConverter[From, To any] interface {
	Convert(from From) (To, error)
}

// BaseGRPCClient - base client for gRPC services.
type BaseGRPCClient struct {
	serviceAddr       string
	conn              *grpc.ClientConn
	DevMode           bool
	ConnectionWarm    bool
	connectionState   int32 // atomic
	mu                sync.RWMutex
	connectCtx        context.Context //nolint:containedctx // it's ok
	connectCancelFunc context.CancelFunc
	connectDone       chan struct{}
	dialOptions       []grpc.DialOption
	connectionTimeout time.Duration
	rpcTimeout        time.Duration
}

func NewBaseGRPCClient[T any](
	serviceAddr string,
	clientCreator ClientCreator[T],
	opts ...ClientOption,
) *GenericGRPCClient[T] {
	ctx, cancel := context.WithCancel(context.Background())

	baseClient := &BaseGRPCClient{
		serviceAddr:       serviceAddr,
		connectionState:   connectionStateDisconnected,
		connectCtx:        ctx,
		connectCancelFunc: cancel,
		connectDone:       make(chan struct{}),
		dialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
		connectionTimeout: defaultConnectionTimeout,
		rpcTimeout:        defaultRpcTimeout,
	}

	for _, opt := range opts {
		opt(baseClient)
	}

	client := &GenericGRPCClient[T]{
		BaseGRPCClient: baseClient,
		clientCreator:  clientCreator,
	}

	go client.connectAsync()

	return client
}

// GenericGRPCClient - generic gRPC client for concrete service.
type GenericGRPCClient[T any] struct {
	*BaseGRPCClient
	client        T
	clientExists  bool
	clientCreator ClientCreator[T]
}

// GetClient returns client for calling gRPC methods.
func (c *GenericGRPCClient[T]) GetClient() (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero T
	if !c.clientExists {
		return zero, ErrServiceNotAvailable
	}

	return c.client, nil
}

// CallRPC runs typed RPC request with error handling.
func (c *GenericGRPCClient[T]) CallRPC(
	ctx context.Context,
	rpcFunc func(T, context.Context) (any, error),
) (any, error) {
	c.mu.RLock()
	client := c.client
	rpcTimeout := c.rpcTimeout
	clientExists := c.clientExists
	c.mu.RUnlock()

	if !clientExists {
		if c.DevMode {
			return nil, nil //nolint:nilnil // ok for dev mode
		}

		return nil, ErrServiceNotAvailable
	}

	if rpcTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, rpcTimeout)
		defer cancel()
	}

	resp, err := rpcFunc(client, ctx)
	if err == nil {
		return resp, nil
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, ErrRPCTimeout
	}

	if status.Code(err) == codes.Unavailable {
		return nil, ErrServiceNotAvailable
	}

	// Ignore error in dev mode
	if c.DevMode {
		return nil, nil //nolint:nilnil // ok for dev mode
	}

	return nil, err
}

// SafeCall is a wrapper for executing typed request with default value
// Additional functions for typed calls should be defined outside
// and used with package functions, because golang methods cannot have
// type params, different from struct params type
// (in short, go struct method cannot have generic params, that not defined in struct)

// ExecuteCall executes typed request for concrete client.
func ExecuteCall[C, ReqT, RespT any](
	client *GenericGRPCClient[C],
	ctx context.Context,
	rpcFunc func(C, ReqT) (RespT, error),
	req ReqT,
) (RespT, error) {
	var zero RespT

	rpcClient, err := client.GetClient()
	if err != nil {
		if client.DevMode {
			return zero, nil
		}

		return zero, err
	}

	client.mu.RLock()
	rpcTimeout := client.rpcTimeout
	client.mu.RUnlock()

	if rpcTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, rpcTimeout)
		defer cancel()
	}

	resp, err := rpcFunc(rpcClient, req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return zero, ErrRPCTimeout
		}

		if status.Code(err) == codes.Unavailable {
			return zero, ErrServiceNotAvailable
		}

		if client.DevMode {
			return zero, nil
		}

		return zero, err
	}

	return resp, nil
}

// ExecuteCallWithFallback executes request with fallback in dev-mode.
func ExecuteCallWithFallback[C, ReqT, RespT any](
	client *GenericGRPCClient[C],
	ctx context.Context,
	rpcFunc func(C, ReqT) (RespT, error),
	req ReqT,
	fallback RespT,
) (RespT, error) {
	var zero RespT

	rpcClient, err := client.GetClient()
	if err != nil {
		if client.DevMode {
			return fallback, nil
		}

		return zero, err
	}

	client.mu.RLock()
	rpcTimeout := client.rpcTimeout
	client.mu.RUnlock()

	if rpcTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, rpcTimeout)
		defer cancel()
	}

	resp, err := rpcFunc(rpcClient, req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return zero, ErrRPCTimeout
		}

		if status.Code(err) == codes.Unavailable {
			return zero, ErrServiceNotAvailable
		}

		if client.DevMode {
			return fallback, nil
		}

		return zero, err
	}

	return resp, nil
}

func (c *BaseGRPCClient) WaitForConnection(ctx context.Context) error {
	state := atomic.LoadInt32(&c.connectionState)
	if state == connectionStateConnected {
		return nil
	}

	if state == connectionStateFailed {
		if c.DevMode {
			return nil
		}

		return ErrConnectionFailed
	}

	select {
	case <-c.connectDone:
		state = atomic.LoadInt32(&c.connectionState)
		if state == connectionStateConnected {
			return nil
		}

		if c.DevMode {
			return nil
		}

		return ErrConnectionFailed
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *GenericGRPCClient[T]) connectAsync() {
	defer close(c.connectDone)

	c.clientExists = false

	atomic.StoreInt32(&c.connectionState, connectionStateConnecting)

	_, cancel := context.WithTimeout(c.connectCtx, c.connectionTimeout)
	defer cancel()

	conn, err := grpc.NewClient(c.serviceAddr, c.dialOptions...)

	if c.connectCtx.Err() != nil {
		atomic.StoreInt32(&c.connectionState, connectionStateFailed)

		return
	}

	if err != nil {
		atomic.StoreInt32(&c.connectionState, connectionStateFailed)

		return
	}

	client := c.clientCreator(conn)

	c.mu.Lock()
	c.conn = conn
	c.client = client
	c.clientExists = true
	c.mu.Unlock()

	// Connection cannot be warmed, because we use generic connections
	// for concrete realisations connection should be warmed independently

	atomic.StoreInt32(&c.connectionState, connectionStateConnected)
}

func (c *BaseGRPCClient) Close() error {
	c.connectCancelFunc()

	<-c.connectDone

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
