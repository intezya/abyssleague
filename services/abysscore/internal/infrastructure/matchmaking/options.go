package matchmaking

type EngineOption func(*engine)

func WithPoolSize(size int) EngineOption {
	return func(e *engine) {
		e.pool = make([]*Entity, size)
	}
}

func WithMaxPointDifference(value uint) EngineOption {
	return func(e *engine) {
		e.maxPointDifferenceForPlayers = value
	}
}
