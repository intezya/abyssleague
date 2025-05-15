package types

import "context"

type Runnable interface {
	Run(ctx context.Context)
}

type RunnableNoCtx interface {
	Run()
}

type Provider[T interface{}] interface {
	Provide() T
}

type Producer[T interface{}] interface {
	Produce() T
}

type Consumer[T interface{}] interface {
	Consume(T)
}
